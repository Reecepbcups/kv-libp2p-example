package main

import (
	"context"
	"fmt"
	"sync"

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
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
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
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("redis command")
		// },
	}
	redisCmd.PersistentFlags().StringP(FlagPeerAddress, "p", "", "peer address")

	// add commands to cmd
	redisCmd.AddCommand(&cobra.Command{
		Use:     "get",
		Short:   "get a value from the redis store",
		Example: `redis get users name -p /ip4/127.0.0.1/tcp/38733/p2p/12D3KooWGEeb4NYtpFwhc7WxQuPGzTf3RvyXKspSMyCrkb5THBzS`,
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get command")
			// node := redis.NewNode()

			// get peer address from flag
			peerAddr, err := cmd.Flags().GetString(FlagPeerAddress)
			if err != nil {
				panic(err)
			}

			store := redis.NewStore("client")

			n := redis.NewServer(store) // is this required for local instances? Should not be since qwe request upstream

			p := getPeer(peerAddr)
			if err := n.Node.Connect(context.Background(), *p); err != nil {
				panic(err)
			}
			fmt.Println("sending req to", peerAddr)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			fmt.Println("sending request to peer before exec")

			respch := n.RedisService.RedisExec(ctx, "users", "name", p.ID)
			res := <-respch
			fmt.Println("Redis Server Response:", res)

		},
	})

	// set
	redisCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "set a value in the redis store",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("set command")
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

			// set some testing values for now
			table := store.Table("users")
			table.Set("name", "John")

			s := redis.NewServer(store)

			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func(s *redis.Server) {
				defer wg.Done()
				s.Stop()
			}(s)

			s.Start()
			wg.Wait()
		},
	}
}
