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

func newRequest(method string, path string, protocol string, host string, userAgent string) *Request {
	return &Request{
		Method:    method,
		Path:      path,
		Protocol:  protocol,
		Host:      host,
		UserAgent: userAgent,
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

	// TODO Add validation
	reqStringSlice := strings.Split(string(buf), "\r\n")
	startLineSlice := strings.Split(reqStringSlice[0], " ")
	host := strings.Split(reqStringSlice[1], ":")[1]
	userAgent := strings.Trim(strings.Split(reqStringSlice[2], ":")[1], " ")
	request := newRequest(startLineSlice[0], startLineSlice[1], startLineSlice[2], host, userAgent)

	fmt.Printf("New Request: %s %s %s\n", request.Method, request.Path, request.Protocol)

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
	response := OK + CRLF + ContentType + CRLF + getContentLen(request.UserAgent) + CRLF + CRLF + request.UserAgent
	conn.Write([]byte(response))

}
func getContentLen(s string) string {
	return strings.Replace(ContentLength, "0", strconv.Itoa(len(s)), 1)
}
