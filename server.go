package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
)

type Server struct {
	Node        host.Host
	PingService *ping.PingService
}

func NewServer() *Server {
	node := NewNode()

	// override the default ping service
	ps := &ping.PingService{
		Host: node,
	}
	node.SetStreamHandler(ping.ID, ps.PingHandler)

	return &Server{
		Node:        node,
		PingService: ps,
	}
}

func NewNode() host.Host {
	// start a libp2p node with default settings
	node, err := libp2p.New(
		libp2p.ForceReachabilityPublic(),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		panic(err)
	}
	return node
}

func (s *Server) Start() {
	fmt.Println("Server started")
	s.PrintPeerInfo()
	fmt.Println("Listen addresses:", s.Node.Addrs())
}

func (s *Server) PrintPeerInfo() {
	// print the node's PeerInfo in multiaddr format
	peerInfo := peerstore.AddrInfo{
		ID:    s.Node.ID(),
		Addrs: s.Node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("libp2p node address:", addrs)
}

func (s *Server) Stop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
	if err := s.Node.Close(); err != nil {
		panic(err)
	}
}
