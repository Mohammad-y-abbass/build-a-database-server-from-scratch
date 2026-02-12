package server

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Mohammad-y-abbass/build-a-database-server-from-scratch/internal/lexer"
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

const (
	colorRed   = "\033[31m"
	colorReset = "\033[0m"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("--- New connection from %s ---\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		query := scanner.Text()
		fmt.Printf("Received query: %q\n", query) // Logs what you typed

		// Create lexer from query string
		l := lexer.New(query)

		// Create parser from lexer
		p := parser.New(l)

		// Parse the program
		program := p.ParseProgram()

		// Check for parsing errors
		if len(p.Errors()) > 0 {
			errorMsg := p.GetErrorMessage()
			fmt.Printf("%sParsing error:%s\n%s\n", colorRed, colorReset, errorMsg)
			conn.Write([]byte(colorRed + errorMsg + colorReset + "\n"))
		} else if program == nil || len(program.Statements) == 0 {
			fmt.Printf("%sError: no statements parsed%s\n", colorRed, colorReset)
			conn.Write([]byte(colorRed + "Error: no statements parsed" + colorReset + "\n"))
		} else {
			response := p.FormatAST(program)
			fmt.Printf("Responding:\n%s\n", response)
			conn.Write([]byte(response + "\n"))
		}
	}
	fmt.Printf("--- Connection closed ---\n")
}
