package main

import (
	"bytes"
	"fmt"
	"io"
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
	BadRequest    = "HTTP/1.1 400 Bad Request"
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
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Printf("Error reading: %v", err)
		os.Exit(1)
	}

	buf = bytes.Trim(buf, "\x00")

	reqString := string(buf[:n])
	reqStringSlice := strings.Split(reqString, CRLF)
	if len(reqStringSlice) < 3 {
		sendResponse(conn, BadRequest, "")
		return
	}

	startLineSlice := strings.Split(reqStringSlice[0], " ")
	if len(startLineSlice) < 3 {
		log.Println("Invalid start line in request")
		sendResponse(conn, BadRequest, "")
		return
	}

	host := strings.SplitN(reqStringSlice[1], ": ", 2)[1]
	userAgent := strings.SplitN(reqStringSlice[2], ": ", 2)[1]

	request := newRequest(startLineSlice[0], startLineSlice[1], startLineSlice[2], host, userAgent)

	fmt.Printf("New Request: %s %s %s\n", request.Method, request.Path, request.Protocol)

	switch {
	case request.Path == "/":
		sendResponse(conn, OK, "")
	case strings.HasPrefix(request.Path, "/echo"):
		handleEcho(conn, request)
	case request.Path == "/user-agent":
		handleUserAgent(conn, request)
	default:
		sendResponse(conn, NotFound, "")
	}

}
func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) == 2 {
			headers[strings.ToLower(parts[0])] = parts[1]
		}
	}
	return headers
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

func sendResponse(conn net.Conn, status string, body string) {
	response := fmt.Sprintf("%s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", status, len(body), body)
	_, err := conn.Write([]byte(response))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
