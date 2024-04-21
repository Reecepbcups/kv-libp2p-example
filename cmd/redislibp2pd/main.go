package main

import (
	"context"
	"fmt"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/reecepbcups/redis-libp2p/redis"
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
	cmd.AddCommand(redisCmd())

	return cmd
}

func startServer() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the redis server",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			store := redis.NewStore("server")
			redis.RunServerNode(store)

			<-ctx.Done()
		},
	}
}

func redisCmd() *cobra.Command {
	redisCmd := &cobra.Command{
		Use:   "redis",
		Short: "redis commands",
	}
	redisCmd.PersistentFlags().StringP(FlagPeerAddress, "p", "", "peer address")

	redisCmd.AddCommand(&cobra.Command{
		Use:     "get",
		Example: `redis get users name -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("get;%s;%s", table, key)

			redis.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "set",
		Example: `redis set table key value -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			value := args[2]
			reqCmd := fmt.Sprintf("set;%s;%s,%s", table, key, value)

			redis.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "del",
		Aliases: []string{"delete"},
		Example: `redis delete table key -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("delete;%s;%s", table, key)

			redis.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "keys",
		Example: `redis keys table -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			reqCmd := fmt.Sprintf("keys;%s", table)

			redis.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "values",
		Example: `redis values table -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			table := args[0]
			reqCmd := fmt.Sprintf("values;%s", table)

			redis.RunClientNode(*getPeerFromFlag(cmd), reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: `redis all -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			redis.RunClientNode(*getPeerFromFlag(cmd), "all")

		},
	})

	return redisCmd
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
