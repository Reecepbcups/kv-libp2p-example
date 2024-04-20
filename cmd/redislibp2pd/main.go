package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}

	redisCmd := &cobra.Command{
		Use:   "redis",
		Short: "redis commands",
		// Run: func(cmd *cobra.Command, args []string) {
		// 	fmt.Println("redis command")
		// },
	}

	// add commands to cmd
	redisCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "get a value from the redis store",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("get command")
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

	cmd.AddCommand(redisCmd)

	return cmd
}
