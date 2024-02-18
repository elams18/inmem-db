package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func printResponse(conn net.Conn) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		response := scanner.Text()
		if strings.HasPrefix(response, "-ERR Empty Command") {
			return
		}
		fmt.Println(response)

		// Check if the response indicates end of data (-1)
		if response == "-1" || strings.HasPrefix(response, "-ERR") {
			//fmt.Println("End of response")
			return
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading response:", err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection", err)
		}
	}(conn)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n') // Read input from user
		_, err := fmt.Fprint(writer, cmd)
		if err != nil {
			fmt.Println("Error from server", err)
			return
		} // Send command to server
		_, err = writer.WriteString("\r\n")
		if err != nil {
			fmt.Println("Error: ", err)
			return
		} // Add newline after the command
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			return
		} // Flush writer to send data immediately

		// Print the response
		printResponse(conn)
	}
}
