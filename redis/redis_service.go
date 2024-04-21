package redis

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

const (
	Protocol   = "/redis/1.0.0"
	PacketSize = 1024
)

func CreateNode() host.Host {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	return node
}

func ReadHelloProtocol(s network.Stream, store *Store) (network.Stream, error) {
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return s, err
	}

	connection := s.Conn()

	fmt.Printf("-> Message from '%s': %s", connection.RemotePeer().String(), message)

	// write data to the stream for the return back
	// if _, err = s.Write([]byte(res + "\n")); err != nil {
	if _, err = s.Write(HandleMsg(message, store)); err != nil {
		return s, err
	}

	// TODO: do this before the newline?
	// pad the rest of the buffer with 0s
	if _, err = s.Write(make([]byte, PacketSize-len(message))); err != nil {
		return s, err
	}

	return s, nil
}

func RunServerNode(store *Store) peerstore.AddrInfo {
	fmt.Printf("Creating server node...")
	targetNode := CreateNode()
	PrintNodeInfo(targetNode)

	targetNode.SetStreamHandler(Protocol, func(s network.Stream) {
		fmt.Printf(Protocol + " stream created!\n")
		if _, err := ReadHelloProtocol(s, store); err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	return *host.InfoFromHost(targetNode)
}

func PrintNodeInfo(node host.Host) {
	peerInfo := peerstore.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peerstore.AddrInfoToP2pAddrs(&peerInfo)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nlibp2p node address: NODE=%s\n", addrs[0])
}

func RunClientNode(targetNodeInfo peerstore.AddrInfo, cmd string) {
	fmt.Println("Creating client node...")
	sourceNode := CreateNode()

	sourceNode.Connect(context.Background(), targetNodeInfo)

	stream, err := sourceNode.NewStream(context.Background(), targetNodeInfo.ID, Protocol)
	if err != nil {
		panic(err)
	}

	if !strings.HasSuffix(cmd, "\n") {
		cmd += "\n"
	}

	fmt.Printf("Sending message...\n")
	_, err = stream.Write([]byte(cmd))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message sent to '%s': %s\n", targetNodeInfo.ID.String(), cmd)

	response := make([]byte, PacketSize)
	n, err := stream.Read(response)
	if err != nil && err.Error() != "EOF" {
		// EOF occurs if packet padding is not added
		fmt.Printf("Error reading response: %s\n", err)
		return
	}

	fmt.Printf("Response from '%s': %s\n", targetNodeInfo.ID.String(), response[:n])
}
