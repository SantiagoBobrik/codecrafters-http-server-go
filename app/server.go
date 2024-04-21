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
	Method    string
	Path      string
	Protocol  string
	Host      string
	UserAgent string
}

const (
	CRLF          = "\r\n"
	OK            = "HTTP/1.1 200 OK"
	NotFound      = "HTTP/1.1 404 Not Found"
	InternalError = "HTTP/1.1 500 Internal Server Error"
	ContentType   = "Content-Type: text/plain"
	ContentLength = "Content-Length: 0"
	UserAgent     = "User-Agent: 0"
)

func newRequest(s []string) *Request {
	return &Request{
		Method:    s[0],
		Path:      s[1],
		Protocol:  s[2],
		Host:      s[3],
		UserAgent: s[4],
	}

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
	request := newRequest(startLineSlice)

	fmt.Printf("Request: %s %s %s\n", request.Method, request.Path, request.Protocol)

	switch {
	case request.Path == "/":
		conn.Write([]byte(OK + CRLF + CRLF))
	case strings.HasPrefix(request.Path, "/echo"):
		handleEcho(conn, request)
	case request.Path == "/user-agent":
		handleUserAgent(conn, request)
	default:
		conn.Write([]byte(NotFound + CRLF + CRLF))
	}

}

func main() {
	listener, err := net.Listen("tcp", "lohalhost:4221")

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
	body, found := strings.CutPrefix(request.Path, "/echo/")

	if !found {
		fmt.Println("Failed to parse request")
		handleServerError(conn)
	}
	response := OK + CRLF + ContentType + CRLF + getContentLen(body) + CRLF + CRLF + body

	_, err := conn.Write([]byte(response))

	if err != nil {
		handleServerError(conn)
	}
}

func handleServerError(conn net.Conn) {
	conn.Write([]byte(InternalError + CRLF + CRLF))
}

func handleUserAgent(conn net.Conn, request *Request) {
	conn.Write([]byte(request.UserAgent))

}
func getContentLen(s string) string {
	return strings.Replace(ContentLength, "0", strconv.Itoa(len(s)), 1)
}
