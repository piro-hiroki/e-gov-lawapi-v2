package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverName    = "e-gov-lawapi-v2"
	serverVersion = "0.1.0"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(serverName + ": ")
	log.SetOutput(os.Stderr)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    serverName,
		Version: serverVersion,
	}, nil)

	registerTools(server, newAPIClient())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
