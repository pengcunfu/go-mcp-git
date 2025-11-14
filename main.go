package main

import (
	"context"
	"log"

	"github.com/pengcunfu/go-mcp-git/internal/server"
	"github.com/spf13/cobra"
)

var (
	repository string
	verbose    int
	userName   string
	userEmail  string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "go-mcp-git",
		Short: "MCP Git Server - Git functionality for MCP",
		Long:  "A Model Context Protocol server providing Git repository interaction and automation tools.",
		Run:   runServer,
	}

	rootCmd.Flags().StringVarP(&repository, "repository", "r", "", "Git repository path")
	rootCmd.Flags().CountVarP(&verbose, "verbose", "v", "Verbose output")
	rootCmd.Flags().StringVarP(&userName, "user-name", "u", "", "Git user name for commits")
	rootCmd.Flags().StringVarP(&userEmail, "user-email", "e", "", "Git user email for commits")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runServer(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	
	srv := server.New(repository, verbose, userName, userEmail)
	if err := srv.Serve(ctx); err != nil {
		log.Fatal(err)
	}
}
