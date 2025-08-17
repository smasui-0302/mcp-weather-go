package main

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	impl := &mcp.Implementation{Name: "weather", Version: "v0.1.0"}
	server := mcp.NewServer(impl, nil)

	registerAlerts(server)
	registerForecast(server)

	// stdioで待ち受ける (stdoutは使わず、ログはstderrに出力される)
	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil && !errors.Is(err, io.EOF) {
		log.Fatalf("failed to run server: %v", err)
	}
}
