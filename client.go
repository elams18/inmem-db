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
		fmt.Fprint(writer, cmd)           // Send command to server
		writer.WriteString("\r\n")        // Add newline after the command
		writer.Flush()                    // Flush writer to send data immediately

		// Read the response from the server
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading response:", err)
			return
		}

		// Print the response
		fmt.Print(response)

		// Check if the response indicates end of data (-1)
		if strings.TrimSpace(response) == "-1\r\n" {
			break
		}
	}
}
