package main

import (
	"context"
	"fmt"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/reecepbcups/kv-libp2p-example/kv"
	"github.com/spf13/cobra"
)

const FlagPeerAddress = "peer"

func main() {
	rootCmd := myCobraCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func myCobraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rootcmd",
		Short: "root cmd",
	}

	cmd.AddCommand(startServer())
	cmd.AddCommand(kvStoreCmd())

	return cmd
}

func startServer() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the KV server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			store := kv.NewStore("server")
			kv.RunServerNode(store)

			<-ctx.Done()
		},
	}
}

func kvStoreCmd() *cobra.Command {
	kvCmd := &cobra.Command{
		Use:     "kv",
		Aliases: []string{"redis", "keyvalue", "key-value", "kvstore"},
	}
	kvCmd.PersistentFlags().StringP(FlagPeerAddress, "p", "", "peer address")

	kvCmd.AddCommand(&cobra.Command{
		Use:     "get",
		Example: `kv get users name -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("get;%s;%s", table, key)

			kv.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	kvCmd.AddCommand(&cobra.Command{
		Use:     "set",
		Example: `kv set table key value -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			value := args[2]
			reqCmd := fmt.Sprintf("set;%s;%s,%s", table, key, value)

			kv.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	kvCmd.AddCommand(&cobra.Command{
		Use:     "del",
		Aliases: []string{"delete"},
		Example: `kv delete table key -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("delete;%s;%s", table, key)

			kv.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	kvCmd.AddCommand(&cobra.Command{
		Use:     "keys",
		Example: `kv keys table -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			reqCmd := fmt.Sprintf("keys;%s", table)

			kv.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	kvCmd.AddCommand(&cobra.Command{
		Use:     "values",
		Example: `kv values table -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			reqCmd := fmt.Sprintf("values;%s", table)

			kv.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	kvCmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: `kv all -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			kv.RunClientNode(*getPeerFromFlag(cmd), "all")

		},
	})

	return kvCmd
}

func getPeerFromFlag(cmd *cobra.Command) *peerstore.AddrInfo {
	peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
	if err != nil {
		panic(err)
	}

	return getPeer(peerAddr)
}

func getPeer(address string) *peerstore.AddrInfo {
	addr, err := multiaddr.NewMultiaddr(address)
	if err != nil {
		panic(err)
	}
	peer, err := peerstore.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	return peer
}
