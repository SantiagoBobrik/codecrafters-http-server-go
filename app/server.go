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
	buf := make([]byte, 4096) // Incrementar tamaño si es necesario
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Error reading: %v", err)
		return // Maneja el error sin terminar el programa
	}

	reqString := string(buf[:n])
	reqStringSlice := strings.Split(reqString, CRLF)
	if len(reqStringSlice) < 3 {
		log.Println("Invalid request")
		return
	}

	startLineSlice := strings.Split(reqStringSlice[0], " ")
	if len(startLineSlice) != 3 {
		log.Println("Invalid start line in request")
		return
	}
	host := strings.TrimSpace(strings.Split(reqStringSlice[1], ": ")[1])
	userAgent := strings.TrimSpace(strings.Split(reqStringSlice[2], ": ")[1])

	request := newRequest(startLineSlice[0], startLineSlice[1], startLineSlice[2], host, userAgent)
	fmt.Printf("New Request: %s %s %s\n", request.Method, request.Path, request.Protocol)

	response := ""
	switch {
	case request.Path == "/":
		response = OK + CRLF + ContentLength + CRLF + CRLF
	case strings.HasPrefix(request.Path, "/echo"):
		handleEcho(conn, request)
	case request.Path == "/user-agent":
		handleUserAgent(conn, request)
	default:
		response = NotFound + CRLF + CRLF
	}

	_, err = conn.Write([]byte(response))
	if err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer listener.Close()

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
