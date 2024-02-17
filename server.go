package main

import (
	"fmt"
	"bufio"
	"net"
	"time"
	"strings"
)

type Database struct{
	data map[string]interface{}
	expiry map[string]time.Time
	sortedSet map[string]map[string]float64
}

func NewDatabase() *Database{
	return &Database{
		data: make(map[string]interface{}),
		expiry: make(map[string]time.Time),
		sortedSet: make(map[string]map[string]float64),
	}
}

func (db *Database) handleCommand(cmd string) string{
	parts := strings.Fields(cmd)
	switch strings.ToUpper(parts[0]){
	case "GET":
		return db.get(parts)
	case "SET":
		return db.set(parts)
	default:
		return fmt.Sprintf("-ERR Unknown command '%s'\r\n", parts[0])
	}
}

func (db *Database) get(parts []string) string{
	if len(parts) != 2 {
		return "-ERR wrong number of arguments for 'GET' command\r\n"
	}
	value, ok := db.data[parts[1]]
	if !ok {
		return "$-1\r\n" // Key not found
	}
	return fmt.Sprintf("$%d\r\n%v\r\n", len(fmt.Sprintf("%v", value)), value)
}

func (db *Database) set(parts []string) string {
	if len(parts) != 3 {
		return "-ERR wrong number of arguments for 'SET' command\r\n"
	}
	db.data[parts[1]] = parts[2]
	return "+OK\r\n"
}


func handleConnection(conn net.Conn, db *Database){
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		cmd, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		cmd = strings.TrimSpace(cmd)
		if cmd == "QUIT" {
			return
		}

		response := db.handleCommand(cmd)
		writer.WriteString(response)
		writer.Flush()
	}
}

func main(){
	db:= NewDatabase()

	listener, err:=net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()
	
	for {
		conn, err:= listener.Accept()
		if err != nil{
			fmt.Println("Error accepting connection", err)
			continue
		}
		go handleConnection(conn, db);
	}
}