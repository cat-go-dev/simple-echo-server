package echoserver

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"net"
)

const defaultPort = 8080

type EchoServer struct {
	port   int
	logger *slog.Logger
}

func NewEchoServer(port int, log *slog.Logger) *EchoServer {
	return &EchoServer{
		port:   port,
		logger: log,
	}
}

func (e EchoServer) Start(ctx context.Context) error {
	p := defaultPort
	if e.port != 0 {
		p = e.port
	}
	port := fmt.Sprintf(":%d", p)

	ln, err := net.Listen("tcp", port)
	if err != nil {
		wErr := fmt.Errorf("Server start failed: %w", err)
		e.logger.ErrorContext(ctx, wErr.Error())
		return wErr
	}

	e.logger.With(slog.String("port", port)).InfoContext(ctx, "Server started")

	inChan := make(chan *net.Conn)
	defer close(inChan)

	go func() {
		for {
			conn, err := ln.Accept()
			l := e.logger.With(slog.String("Client remote address", conn.RemoteAddr().String()))
			if err != nil {
				l.ErrorContext(ctx, fmt.Sprintf("Connection accepting failed: %v", err))
				continue
			}

			l.InfoContext(ctx, "Client connected")

			go func() {
				inChan <- &conn
			}()
		}
	}()

	for {
		select {
		case conn := <-inChan:
			if err := e.handle(ctx, conn); err != nil {
				c := *conn
				e.logger.With(slog.String("Client remote address", c.RemoteAddr().String())).ErrorContext(ctx, fmt.Sprintf("Handling message failed: %v", err))
			}
		case <-ctx.Done():
			e.logger.InfoContext(ctx, "Context done")
			return nil
		}
	}
}

func (e EchoServer) handle(ctx context.Context, conn *net.Conn) error {
	c := *conn
	defer c.Close()

	reader := bufio.NewReader(c)
	msg, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	e.logger.With(slog.String("Client remote address", c.RemoteAddr().String())).InfoContext(ctx, fmt.Sprintf("Message received: %s", msg))

	writter := bufio.NewWriter(c)
	if _, err := writter.WriteString(fmt.Sprintf("Get your message back: %s", msg)); err != nil {
		return err
	}
	if err := writter.Flush(); err != nil {
		return err
	}

	e.logger.With(slog.String("Client remote address", c.RemoteAddr().String())).InfoContext(ctx, "Message sent")

	return nil
}
