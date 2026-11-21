package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Server struct {
	handler *commands.Handler
}

func New(handler *commands.Handler) *Server {
	return &Server{
		handler: handler,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind to port 6379: %w", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("error accepting connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			log.Printf("error reading from connection: %v", err)
			break
		}

		args, err := resp.Parse(buf[:n])
		if err != nil || len(args) == 0 {
			log.Printf("error parsing input: %v", err)
			return
		}

		resp, err := s.handler.Handle(args)
		if err != nil {
			log.Print(err)
			continue
		}

		if err = writeToConn(conn, resp); err != nil {
			log.Print(err)
			break
		}
	}
}

func writeToConn(conn net.Conn, data string) error {
	_, err := conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error writing a response: %w", err)
	}

	return nil
}
