package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Request struct {
	method   string
	path     string
	protocol string
}

var successResponse = "HTTP/1.1 200 OK\r\n\r\n"
var notFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
var contentType = "Content-Type: text/plain \r\n\r\n"

func newRequest(method string, path string, protocol string) *Request {
	return &Request{method, path, protocol}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	if err != nil {
		log.Fatal("Failed to read data")
	}

	reqStringSlice := strings.Split(string(buf), "\r\n")
	startLineSlice := strings.Split(reqStringSlice[0], " ")
	request := newRequest(startLineSlice[0], startLineSlice[1], startLineSlice[2])
	ok := false

	fmt.Printf("Request: %s %s %s\n", request.method, request.path, request.protocol)

	handleRoute("/echo", request, func(r *Request) {
		ok = true
		paths := strings.Split(request.path, "/")
		lastPath := paths[len(paths)-1]
		conn.Write([]byte(successResponse + contentType + lastPath))
	})

	if !ok {
		conn.Write([]byte(notFoundResponse))
	}

}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection")
		}
		go handleConnection(conn)
	}
}

func handleRoute(path string, request *Request, handler func(r *Request)) {
	hasPrefix := strings.HasPrefix(request.path, path)

	if hasPrefix {
		handler(request)
	}

}
