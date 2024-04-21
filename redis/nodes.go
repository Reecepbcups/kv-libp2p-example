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

const (
	Protocol   = "/hello/1.0.0"
	PacketSize = 1024
)

func CreateNode() host.Host {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	return node
}

func ReadHelloProtocol(s network.Stream) (network.Stream, error) {
	// TO BE IMPLEMENTED: Read the stream and print its content
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return s, err
	}

	connection := s.Conn()

	fmt.Printf("-> Message from '%s': %s", connection.RemotePeer().String(), message)

	// write data to the stream for the return back
	_, err = s.Write([]byte("Hello from the other side!\n"))
	if err != nil {
		return s, err
	}
	// pad the rest of the buffer with 0s
	_, err = s.Write(make([]byte, PacketSize-len(message)))
	if err != nil {
		return s, err
	}

	// return nil
	return s, nil
}

// Targert = server
func RunTargetNode() peerstore.AddrInfo {
	fmt.Printf("Creating target node...")
	targetNode := CreateNode()
	PrintNodeInfo(targetNode)

	// TO BE IMPLEMENTED: Set stream handler for the "/hello/1.0.0" protocol
	targetNode.SetStreamHandler(Protocol, func(s network.Stream) {
		fmt.Printf(Protocol + " stream created!\n")
		if _, err := ReadHelloProtocol(s); err != nil {
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
	stream, err := sourceNode.NewStream(context.Background(), targetNodeInfo.ID, Protocol)
	if err != nil {
		panic(err)
	}

	message := "Hello from Launchpad!\n"
	fmt.Printf("Sending message...\n")
	_, err = stream.Write([]byte(message))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Message sent to '%s': %s\n", targetNodeInfo.ID.String(), message)

	// resp := make(chan string)

	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// go func() {
	// 	defer close(resp)

	// 	fmt.Println("Reading response...")
	// 	fmt.Println("Response:", <-resp)

	// 	stream.Close()
	// }()

	newS, err := ReadHelloProtocol(stream)
	if err != nil {
		stream.Reset()
	}
	// wg.Wait()

	response := make([]byte, PacketSize)
	n, err := newS.Read(response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response from '%s': %s\n", targetNodeInfo.ID.String(), response[:n])

	// if err = ReadHelloProtocol(stream, resp); err != nil {
	// 	stream.Reset()
	// }

	// stream.Close()
}
