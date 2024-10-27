package main

import (
	"context"
	"log/slog"
	echoserver "mseaps/echo-server"
	"os"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	echoServer := echoserver.NewEchoServer(5000, logger)

	echoServer.Start(ctx)
}
