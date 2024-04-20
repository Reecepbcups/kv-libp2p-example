package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// iterate over the table

func main() {
	// fmt.Println("Hello, World!")

	// store := NewStore()
	// table := store.Table("users")

	// table.Set("name", "John")
	// table.Set("age", fmt.Sprintf("%d", 30))

	// name, _ := table.Get("name")
	// fmt.Println(name)

	// age, _ := table.Get("age")
	// fmt.Println(age)

	s := NewServer()

	if clientCommand(s) {
		// client only command line stuff, no server here with arg
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(s *Server) {
		defer wg.Done()
		s.Stop()
	}(s)

	s.Start()
	wg.Wait()
}

func clientCommand(s *Server) bool {
	if len(os.Args) <= 1 {
		return false
	}

	fmt.Println("Command line arguments:")
	for i, arg := range os.Args {
		fmt.Printf("arg %d: %s\n", i, arg)
	}

	addr, err := multiaddr.NewMultiaddr(os.Args[1])
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}

	if err := s.Node.Connect(context.Background(), *peer); err != nil {
		panic(err)
	}
	fmt.Println("sending 5 ping messages to", addr)

	ch := s.PingService.Ping(context.Background(), peer.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("got ping response!", "RTT:", res.RTT)
	}

	return true
}
