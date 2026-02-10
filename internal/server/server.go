package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/Mohammad-y-abbass/build-a-database-server-from-scratch/internal/parser"
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("--- New connection from %s ---\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		query := scanner.Text()
		fmt.Printf("Received query: %q\n", query) // Logs what you typed

		p := parser.New(query)
		err := p.Parse()

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			conn.Write([]byte("parsing_error\n"))
		} else {
			response := formatOutput(p.GetOutput())
			fmt.Printf("Responding: %s\n", response)
			conn.Write([]byte(response + "\n"))
		}
	}
	fmt.Printf("--- Connection closed ---\n")
}

func formatOutput(output []any) string {
	if len(output) == 0 {
		return "no rows"
	}
	var res []string
	for _, val := range output {
		switch v := val.(type) {
		case bool:
			res = append(res, strings.ToUpper(fmt.Sprint(v))) // TRUE not true
		default:
			res = append(res, fmt.Sprint(v))
		}
	}
	return strings.Join(res, ", ")
}
