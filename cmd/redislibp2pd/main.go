package main

import (
	"context"
	"fmt"

	peerstore "github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/reecepbcups/redis-libp2p/redis"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := myCobraCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
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

func myCobraCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rootcmd",
		Short: "root cmd",
	}

	cmd.AddCommand(startServer())
	cmd.AddCommand(redisCmd())

	return cmd
}

func redisCmd() *cobra.Command {
	FlagPeerAddress := "peer"

	redisCmd := &cobra.Command{
		Use:   "redis",
		Short: "redis commands",
	}
	redisCmd.PersistentFlags().StringP(FlagPeerAddress, "p", "", "peer address")

	// add commands to cmd
	redisCmd.AddCommand(&cobra.Command{
		Use:     "get",
		Short:   "get a value from the redis store",
		Example: `redis get users name -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get command")
			// node := redis.NewNode()

			// get peer address from flag
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("get;%s;%s", table, key)

			p := getPeer(peerAddr)
			redis.RunClientNode(*p, reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "set",
		Short:   "set a value from the redis store",
		Example: `redis set table key value -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			table := args[0]
			key := args[1]
			value := args[2]
			reqCmd := fmt.Sprintf("set;%s;%s,%s", table, key, value)

			p := getPeer(peerAddr)
			redis.RunClientNode(*p, reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "del",
		Aliases: []string{"delete"},
		Short:   "delete a value from the redis store",
		Example: `redis delete table key -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			table := args[0]
			key := args[1]
			reqCmd := fmt.Sprintf("delete;%s;%s", table, key)

			p := getPeer(peerAddr)
			redis.RunClientNode(*p, reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "keys",
		Short:   "keys from the redis store",
		Example: `redis keys table -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			table := args[0]
			reqCmd := fmt.Sprintf("keys;%s", table)

			p := getPeer(peerAddr)
			redis.RunClientNode(*p, reqCmd)

		},
	})

	redisCmd.AddCommand(&cobra.Command{
		Use:     "all",
		Short:   "all values",
		Example: `redis all -p /ip4/127.0.0.1/tcp/38733/p2p/XXXXXXX`,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			p := getPeer(peerAddr)
			redis.RunClientNode(*p, "all")

		},
	})

	return redisCmd
}

func startServer() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the redis server",
		Run: func(cmd *cobra.Command, args []string) {
			store := redis.NewStore("server")

			// s := redis.NewServer(store)

			// wg := &sync.WaitGroup{}
			// wg.Add(1)
			// go func(s *redis.Server) {
			// 	defer wg.Done()
			// 	s.Stop()
			// }(s)

			// s.Start()
			// wg.Wait()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			redis.RunServerNode(store)
			// redis.RunSourceNode(info)

			<-ctx.Done()
		},
	}
}
