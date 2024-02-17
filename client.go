package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n') // Read input from user
		cmd = strings.TrimSpace(cmd)

		// Send command to server
		fmt.Fprintln(writer, cmd)
		writer.Flush()

		// Read and print response from the server
		printResponse(conn)
	}
}

func printResponse(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// Trim whitespace and newlines
		response = strings.TrimSpace(response)

		// Check if the response indicates end of data (-1)
		if response == "-1" {
			fmt.Println("End of response")
			return
		}

		// Check if the response has a prefix character indicating response type
		if len(response) >= 1 {
			switch response[0] {
			case '+', '-', ':':
				// Single-line response, print and exit loop
				fmt.Println(response)
				return
			case '$':
				// Handle quoted strings for bulk responses
				if strings.HasPrefix(response, "$\"") && strings.HasSuffix(response, "\"") {
					// Extract the quoted string value
					value := response[2 : len(response)-1]
					fmt.Println(value)
				} else {
					fmt.Println("Unexpected response format:", response)
				}
				return
			case '*':
				// Multi-line response, print and continue loop
				fmt.Println(response)
				for {
					// Read additional lines of multi-line response
					line, err := reader.ReadString('\n')
					if err != nil {
						fmt.Println("Error reading additional response line:", err)
						return
					}
					line = strings.TrimSpace(line)
					fmt.Println(line)
					if line == "-1" {
						return // End of multi-line response
					}
				}
			default:
				// Unexpected response type, print and exit loop
				fmt.Println("Unexpected response:", response)
				return
			}
		} else {
			// Unexpected response format, print and exit loop
			fmt.Println("Unexpected response format:", response)
			return
		}
	}
}
