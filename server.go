package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

type Server struct {
	node host.Host // set on start
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	fmt.Println("Server started")

	// start a libp2p node with default settings
	node, err := libp2p.New(libp2p.ForceReachabilityPublic(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	s.node = node

	// print the node's listening addresses
	fmt.Println("Listen addresses:", node.Addrs())
}

func (s *Server) Stop() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	fmt.Println("Received signal, shutting down...")
}
