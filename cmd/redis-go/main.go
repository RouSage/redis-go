package main

import (
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func main() {
	log.SetFlags(0)

	listener, err := net.Listen("tcp", "0.0.0.0:6379") //nolint:gosec // G102: must listen on all interfaces for the tester
	if err != nil {
		log.Fatalf("failed to bind to port 6379: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from connection: %v", err)
			break
		}

		args, err := resp.Parse(buf[:n])
		if err != nil {
			log.Fatalf("error parsing input: %v", err)
		}

		switch strings.ToUpper(args[0]) {
		case "PING":
			writeToConn(conn, resp.SimpleString("PONG"))
			return
		case "ECHO":
			if len(args) < 2 {
				log.Printf("ECHO call needs %d arg but got %d args", 2, len(args))
			}
			writeToConn(conn, resp.BulkString(args[1]))
			return
		}
	}
}

func writeToConn(conn net.Conn, data string) {
	_, err := conn.Write([]byte(data))
	if err != nil {
		log.Printf("error writing a response: %v", err)
		return
	}
}
