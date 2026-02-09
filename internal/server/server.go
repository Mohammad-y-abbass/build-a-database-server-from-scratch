package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func Start() {
	listener, err := net.Listen("tcp", ":3003")
	if err != nil {
		fmt.Println("Could not start server: ", err)
	}
	defer listener.Close()

	fmt.Println("Server is running on port 3003")

	for {
		// Wait for a client to connect
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected: ", conn.RemoteAddr())

	// Set up a buffered reader to look for the null byte
	reader := bufio.NewReader(conn)

	for {
		// 1. Read until the null byte \x00
		message, err := reader.ReadBytes('\x00')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error: ", err)
			}
			break // Connection closed by client
		}

		// 2. Log what we received (useful for debugging)
		fmt.Println("Received: ", string(message))

		// 3. Respond back with "hello" followed by a null byte
		// The test suite requires the null byte to know the response is finished
		_, err = conn.Write([]byte("hello\x00"))
		if err != nil {
			fmt.Println("Write error: ", err)
			break
		}
	}
	fmt.Println("Client disconnected")
}
