package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	logging "github.com/ipfs/go-log/v2"
	pool "github.com/libp2p/go-buffer-pool"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

var log = logging.Logger("redis")

const (
	// reqTimeout = time.Second * 60

	ID = "/ipfs/redis/1.0.0"

	ServiceName = "libp2p.redis"
)

type RedisService struct {
	Host       host.Host
	RedisStore *Store
}

func NewRedisService(h host.Host, r *Store) *RedisService {
	rs := &RedisService{
		Host:       h,
		RedisStore: r,
	}
	h.SetStreamHandler(ID, rs.RedisHandler)
	return rs
}

func (p *RedisService) RedisHandler(s network.Stream) {
	if err := s.Scope().SetService(ServiceName); err != nil {
		log.Debugf("error attaching stream to ping service: %s", err)
		s.Reset()
	}

	fmt.Println("Stream details:", s)
}

// Client side
type Result struct {
	Resp  any
	Error error
}

func (rs *RedisService) RedisExec(ctx context.Context, table, key string, p peer.ID) <-chan Result {
	return RedisExecute(ctx, rs.Host, rs.RedisStore, table, key, p)
}

func redisError(err error) chan Result {
	ch := make(chan Result, 1)
	ch <- Result{Error: err}
	close(ch)
	return ch
}

func RedisExecute(ctx context.Context, h host.Host, store *Store, table, key string, p peer.ID) <-chan Result {
	s, err := h.NewStream(network.WithUseTransient(ctx, "ping"), p, ID)
	if err != nil {
		return redisError(err)
	}

	if err := s.Scope().SetService(ServiceName); err != nil {
		log.Debugf("error attaching stream to ping service: %s", err)
		s.Reset()
		return redisError(err)
	}

	ctx, cancel := context.WithCancel(ctx)

	out := make(chan Result)
	go func() {
		defer close(out)
		defer cancel()

		for ctx.Err() == nil {
			var res Result
			res.Resp, res.Error = redis(s, store, table, key)

			// canceled, ignore everything.
			if ctx.Err() != nil {
				return
			}

			// No error, record the RTT.
			// if res.Error == nil {
			// 	h.Peerstore().RecordLatency(p, res.Resp)
			// }

			select {
			case out <- res:
			case <-ctx.Done():
				return
			}
		}
	}()
	go func() {
		// forces the ping to abort.
		<-ctx.Done()
		s.Reset()
	}()

	return out
}

func redis(s network.Stream, rstore *Store, table, key string) (string, error) {
	value, ok := rstore.Table(table).Get(key)
	if !ok {
		return "", errors.New(fmt.Sprintf("key %s not found in table %s", key, table))
	}

	size := bytes.NewBufferString(value).Len()

	if err := s.Scope().ReserveMemory(2*size, network.ReservationPriorityAlways); err != nil {
		log.Debugf("error reserving memory for ping stream: %s", err)
		s.Reset()
		return "", err
	}
	defer s.Scope().ReleaseMemory(2 * size) // idk if I need this?

	buf := pool.Get(size)
	defer pool.Put(buf)

	if _, err := io.ReadFull(bytes.NewBufferString(value), buf); err != nil {
		return "", err
	}

	if _, err := s.Write(buf); err != nil {
		return "", err
	}

	return value, nil
}
