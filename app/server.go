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

var request = &Request{}
var successResponse = "HTTP/1.1 200 OK\r\n\r\n"
var contentType = "Content-Type: text/plain \r\n\r\n"

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

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	_, err := conn.Read(buf)

	if err != nil {
		log.Fatal("Failed to read data")
	}

	reqStringSlice := strings.Split(string(buf), "\r\n")
	startLineSlice := strings.Split(reqStringSlice[0], " ")

	request.method = startLineSlice[0]
	request.path = startLineSlice[1]
	request.protocol = startLineSlice[2]

	fmt.Printf("Request: %s %s %s\n", request.method, request.path, request.protocol)

	params := strings.Split(request.path, "/")

	content := params[len(params)-1] + "\r\n\r\n"

	conn.Write([]byte(successResponse + contentType + content))

}
