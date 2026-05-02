package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/piro-hiroki/e-gov-lawapi-v2/internal/egov"
	"github.com/piro-hiroki/e-gov-lawapi-v2/internal/mcpserver"
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

	mcpserver.RegisterTools(server, egov.NewClient(&egov.Options{
		UserAgent: serverName + "/" + serverVersion,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
