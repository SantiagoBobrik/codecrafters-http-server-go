//package pepe
//
//import (
//	"bytes"
//	"fmt"
//	"io"
//	"net"
//	"os"
//	"strings"
//	// Uncomment this block to pass the first stage
//	// "net"
//	// "os"
//)
//
//func main() {
//	// You can use print statements as follows for debugging, they'll be visible when running tests.
//	fmt.Println("Logs from your program will appear here!")
//	// Uncomment this block to pass the first stage
//	l, err := net.Listen("tcp", "0.0.0.0:4221")
//	if err != nil {
//		fmt.Println("Failed to bind to port 4221")
//		os.Exit(1)
//	}
//	for {
//		con, err := l.Accept()
//		if err != nil {
//			fmt.Println("Error accepting connection: ", err.Error())
//			os.Exit(1)
//		}
//
//		go handleConn(con)
//	}
//}
//func handleConn(conn net.Conn) {
//	defer func(conn net.Conn) {
//		err := conn.Close()
//		if err != nil {
//			exitWithError(err)
//		}
//	}(conn)
//	clientMsg := make([]byte, 1024)
//	_, err := conn.Read(clientMsg)
//	if err != nil && err != io.EOF {
//		exitWithError(err)
//	}
//	clientMsg = bytes.Trim(clientMsg, "\x00")
//	headers := map[string]string{}
//	lines := strings.Split(string(clientMsg), "\r\n")
//	if len(lines) > 0 {
//		// Parse headers
//		for _, line := range lines[1:] {
//			if line == "" {
//				break
//			}
//			header := strings.Split(line, ": ")
//			headers[strings.ToLower(header[0])] = header[1]
//		}
//		lineSplit := strings.Split(lines[0], " ")
//		if len(lineSplit) >= 2 {
//			switch lineSplit[1] {
//			case "/":
//				sendResponse(conn, []byte("HTTP/1.1 200 OK\r\n\r\n"))
//			case "/user-agent":
//				if userArgent, ok := headers["user-agent"]; ok {
//					sendResponse(conn, []byte(
//						fmt.Sprintf(
//							"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userArgent), userArgent,
//						)),
//					)
//				}
//				sendResponse(conn, []byte(
//					fmt.Sprintf(
//						"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", 0, "",
//					)),
//				)
//			default:
//				if strings.HasPrefix(lineSplit[1], "/echo/") {
//					respBody := strings.TrimPrefix(lineSplit[1], "/echo/")
//					sendResponse(conn, []byte(
//						fmt.Sprintf(
//							"HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(respBody), respBody,
//						)),
//					)
//					break
//				}
//				sendResponse(conn, []byte("HTTP/1.1 404 Not Found\r\n\r\n"))
//			}
//		}
//	}
//}
//func exitWithError(err error) {
//	fmt.Println(err)
//	os.Exit(1)
//}
//func sendResponse(conn net.Conn, data []byte) {
//	_, err := conn.Write(data)
//	if err != nil {
//		exitWithError(err)
//	}
//}
