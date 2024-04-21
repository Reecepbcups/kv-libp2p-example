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

func HandleMsg(msg string, store *Store) string {
	// TODO: add protocol version here? (smart but not necessary for example demo)
	// formats:
	// 1: set;table;key,value
	// 2: get;table;key

	msg = strings.TrimSuffix(msg, "\n")

	args := strings.Split(msg, ";")
	fmt.Println(args)

	action := args[0]
	table := args[1]

	switch action {
	case "set":
		tuple := strings.Split(args[2], ",")
		key, value := tuple[0], tuple[1]

		store.Table(table).Set(key, value)
		return "OK"
	case "get":
		key := args[2]
		res, ok := store.Table(table).Get(key)
		if !ok {
			return fmt.Sprintf("Key '%s' not found in table '%s'", key, table)
		}
		return res
	default:
		fmt.Println("Invalid message format")
	}

	return ""
}

func ReadHelloProtocol(s network.Stream, store *Store) (network.Stream, error) {
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return s, err
	}

	connection := s.Conn()

	fmt.Printf("-> Message from '%s': %s", connection.RemotePeer().String(), message)

	res := HandleMsg(message, store)

	// write data to the stream for the return back
	_, err = s.Write([]byte(res + "\n"))
	if err != nil {
		return s, err
	}

	// pad the rest of the buffer with 0s (TODO: do this before the newline?)
	_, err = s.Write(make([]byte, PacketSize-len(message)))
	if err != nil {
		return s, err
	}

	return s, nil
}

func RunServerNode(store *Store) peerstore.AddrInfo {
	fmt.Printf("Creating target node...")
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
	fmt.Println("libp2p node address:", addrs[0])
}

func RunClientNode(targetNodeInfo peerstore.AddrInfo, cmd string) {
	fmt.Printf("Creating source node...")
	sourceNode := CreateNode()

	sourceNode.Connect(context.Background(), targetNodeInfo)

	// TO BE IMPLEMENTED: Open stream and send message
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
	// we can ignore padding if we just throw out panic: EOF replies
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response from '%s': %s\n", targetNodeInfo.ID.String(), response[:n])
}
