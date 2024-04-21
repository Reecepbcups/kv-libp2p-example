package redis

import (
	"bufio"
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
)

func CreateNode() host.Host {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	return node
}

func ReadHelloProtocol(s network.Stream) error {
	// TO BE IMPLEMENTED: Read the stream and print its content
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	connection := s.Conn()

	fmt.Printf("Message from '%s': %s", connection.RemotePeer().String(), message)
	return nil
}

// Targert = server
func RunTargetNode() peerstore.AddrInfo {
	fmt.Printf("Creating target node...")
	targetNode := CreateNode()
	PrintNodeInfo(targetNode)

	// TO BE IMPLEMENTED: Set stream handler for the "/hello/1.0.0" protocol
	targetNode.SetStreamHandler("/hello/1.0.0", func(s network.Stream) {
		fmt.Printf("/hello/1.0.0 stream created")
		err := ReadHelloProtocol(s)
		if err != nil {
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
	fmt.Println("libp2p node address:", addrs[0])
}

func RunSourceNode(targetNodeInfo peerstore.AddrInfo) {
	fmt.Printf("Creating source node...")
	sourceNode := CreateNode()
	// fmt.Printf("Source node created with ID '%s'", sourceNode.ID().String())

	sourceNode.Connect(context.Background(), targetNodeInfo)

	// TO BE IMPLEMENTED: Open stream and send message
	stream, err := sourceNode.NewStream(context.Background(), targetNodeInfo.ID, "/hello/1.0.0")
	if err != nil {
		panic(err)
	}

	message := "Hello from Launchpad!\n"
	fmt.Printf("Sending message...")
	_, err = stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	// print resp
	err = ReadHelloProtocol(stream)
	if err != nil {
		stream.Reset()
	} else {
		stream.Close()
	}

	fmt.Printf("Message sent to '%s': %s", targetNodeInfo.ID.String(), message)
}
