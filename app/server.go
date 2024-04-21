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

const (
	CRLF           = "\r\n"
	OK             = "HTTP/1.1 200 OK"
	NOT_FOUND      = "HTTP/1.1 404 Not Found"
	INTERNAL_ERROR = "HTTP/1.1 500 Internal Server Error"
	CONTENT_TYPE   = "Content-Type: text/plain"
	CONTENT_LENGTH = "Content-Length: 0"
)

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
		conn.Write([]byte(OK + CRLF + CRLF))
	case strings.HasPrefix(request.path, "/echo"):
		handleEcho(conn, request)
	default:
		conn.Write([]byte(NOT_FOUND))
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
	body, found := strings.CutPrefix(request.path, "/echo/")

	if !found {
		fmt.Println("Failed to parse request")
		handleServerError(conn)
	}
	response := OK + CRLF + CONTENT_TYPE + CRLF + getContentLen(body) + CRLF + CRLF + body

	_, err := conn.Write([]byte(response))

	if err != nil {
		handleServerError(conn)
	}
}

func handleServerError(conn net.Conn) {
	conn.Write([]byte(INTERNAL_ERROR))
}

func getContentLen(s string) string {
	return strings.Replace(CONTENT_LENGTH, "0", strconv.Itoa(len(s)), 1)
}
