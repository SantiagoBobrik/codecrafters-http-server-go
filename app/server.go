package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection from ", conn.RemoteAddr())

	_, err := conn.Write([]byte("Hello World!"))

	if err != nil {
		return
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	_, err = l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)

	}

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Fatalln("Error accepting connection: ", err.Error())
		}

		go handleConnection(conn)
	}

}
