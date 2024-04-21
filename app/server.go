package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Request struct {
	method   string
	path     string
	protocol string
}

var successResponse = "HTTP/1.1 200 OK\r\n\r\n"
var notFoundResponse = "HTTP/1.1 404 Not Found\r\n\r\n"
var serverErrorResponse = "HTTP/1.1 500 Internal Server Error\r\n\r\n"
var contentType = "Content-Type: text/plain \r\n\r\n"
var contentLength = "Content-Length: 0\r\n\r\n"

func newRequest(method string, path string, protocol string) *Request {
	return &Request{method, path, protocol}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	if err != nil {
		handleServerError(conn)
		log.Fatal("Failed to read data")
	}

	reqStringSlice := strings.Split(string(buf), "\r\n")
	startLineSlice := strings.Split(reqStringSlice[0], " ")
	request := newRequest(startLineSlice[0], startLineSlice[1], startLineSlice[2])

	fmt.Printf("Request: %s %s %s\n", request.method, request.path, request.protocol)

	switch {
	case request.path == "/":
		conn.Write([]byte(successResponse))
	case strings.HasPrefix(request.path, "/echo"):
		handleEcho(conn, request)
	default:
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

func handleEcho(conn net.Conn, request *Request) {
	paths := strings.Split(request.path, "/")
	lastPath := paths[len(paths)-1]

	_, err := conn.Write([]byte(successResponse + contentType + getContentLen(lastPath) + lastPath))
	if err != nil {
		handleServerError(conn)
	}
}

func handleServerError(conn net.Conn) {
	conn.Write([]byte(serverErrorResponse))
}

func getContentLen(s string) string {
	return strings.Replace(contentLength, "0", strconv.Itoa(len(s)), 1)
}
