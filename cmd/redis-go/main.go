package main

import (
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to bind to port 6379: %w", err))
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error accepting connection: %w", err))
		os.Exit(1)
	}

	_, err = conn.Write([]byte(resp.String("PONG")))
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error writing a response: %w", err))
		os.Exit(1)
	}
}
