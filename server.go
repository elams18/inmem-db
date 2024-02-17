// server.go
package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Database struct {
	data      map[string]string
	expiry    map[string]time.Time
	sortedSet map[string]map[string]float64
}

func NewDatabase() *Database {
	return &Database{
		data:      make(map[string]string),
		expiry:    make(map[string]time.Time),
		sortedSet: make(map[string]map[string]float64),
	}
}

func (db *Database) handleCommand(command string) string {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "-ERR Empty command\r\n"
	}

	switch strings.ToUpper(parts[0]) {
	case "GET":
		return db.get(parts)
	case "SET":
		return db.set(splitCommand(command))
	case "DEL":
		return db.del(parts)
	case "EXPIRE":
		return db.expire(parts)
	case "KEYS":
		if len(parts) != 2 {
			return "-ERR wrong number of arguments for 'KEYS' command\r\n"
		}
		keys := db.keys(parts[1])
		if keys == "-1\r\n" {
			return keys // No key matches
		}
		return keys // Add single \r\n after entire response is constructed

	case "TTL":
		return db.ttl(parts)
	case "ZADD":
		return db.zadd(parts)
	case "ZRANGE":
		return db.zrange(parts)
	default:
		return fmt.Sprintf("-ERR Unknown command '%s'\r\n", parts[0])
	}
}

func splitCommand(cmd string) []string {
	var parts []string
	inQuotes := false
	currentPart := ""
	for _, char := range cmd {
		if char == '"' {
			inQuotes = !inQuotes
		} else if char == ' ' && !inQuotes {
			if currentPart != "" {
				parts = append(parts, currentPart)
				currentPart = ""
			}
			continue
		}
		currentPart += string(char)
	}
	if currentPart != "" {
		parts = append(parts, currentPart)
	}
	return parts
}

func (db *Database) get(parts []string) string {
	if len(parts) != 2 {
		return "-ERR wrong number of arguments for 'GET' command\r\n"
	}
	value, ok := db.data[parts[1]]
	if !ok {
		return "$-1\r\n" // Key not found
	}
	key := parts[1]
	expireTime, ok := db.expiry[key]
	if ok && expireTime.Before(time.Now()) {
		delete(db.data, key)   // Remove the expired key
		delete(db.expiry, key) // Remove the expiration time entry
		return "$-1\r\n"       // Key has expired
	}
	return fmt.Sprintf("%s\r\n", value)
}

func (db *Database) set(parts []string) string {
	if len(parts) != 3 && len(parts) != 5 {
		return "-ERR wrong number of arguments for 'SET' command\r\n"
	}
	key := parts[1]
	value := parts[2]
	db.data[key] = value
	if len(parts) >= 5 && strings.ToUpper(parts[3]) == "EX" {
		expireTime, err := strconv.Atoi(parts[4])
		if err != nil {
			return "-ERR invalid expiration time\r\n"
		}
		db.expiry[key] = time.Now().Add(time.Second * time.Duration(expireTime))
		// Set expiration time using a goroutine
		go func(key string, expireTime int) {
			<-time.After(time.Duration(expireTime) * time.Second)
			delete(db.data, key)
		}(key, expireTime)
	}
	return "+OK\r\n"
}

func (db *Database) del(parts []string) string {
	if len(parts) < 2 {
		return "-ERR wrong number of arguments for 'DEL' command\r\n"
	}
	count := 0
	for i := 1; i < len(parts); i++ {
		if _, ok := db.data[parts[i]]; ok {
			delete(db.data, parts[i])
			count++
		}
	}
	return fmt.Sprintf(":%d\r\n", count)
}

func (db *Database) expire(parts []string) string {
	if len(parts) != 3 {
		return "-ERR wrong number of arguments for 'EXPIRE' command\r\n"
	}
	key := parts[1]
	expiry, err := strconv.Atoi(parts[2])
	if err != nil {
		return "-ERR invalid expire time\r\n"
	}
	db.expiry[key] = time.Now().Add(time.Second * time.Duration(expiry))
	return "$:1\r\n"
}

func match(pattern, key string) bool {
	i, j := 0, 0
	for i < len(pattern) && j < len(key) {
		if pattern[i] == '?' || pattern[i] == key[j] {
			i++
			j++
		} else if pattern[i] == '*' {
			if i+1 < len(pattern) && pattern[i+1] == key[j] {
				i++
			} else {
				j++
			}
		} else {
			return false
		}
	}
	return i == len(pattern) && j == len(key)
}

func (db *Database) keys(pattern string) string {
	if pattern == "*" {
		var response strings.Builder
		for key := range db.data {
			response.WriteString(fmt.Sprintf("\"%s\"\r\n", key))
		}
		response.WriteString("-1\r\n") // Indicate end of response
		return response.String()
	}

	var result []string
	for key := range db.data {
		if match(pattern, key) {
			result = append(result, key)
		}
	}
	if len(result) == 0 {
		return "-1\r\n" // Return -1\r\n if there are no key matches
	}
	var response strings.Builder
	for _, key := range result {
		response.WriteString(fmt.Sprintf("\"%s\"\r\n", key))
	}
	response.WriteString("-1\r\n") // Indicate end of response
	return response.String()
}

func (db *Database) ttl(parts []string) string {
	if len(parts) != 2 {
		return "-ERR wrong number of arguments for 'TTL' command\r\n"
	}
	key := parts[1]
	if expiry, ok := db.expiry[key]; ok {
		ttl := expiry.Sub(time.Now())
		if ttl > 0 {
			return fmt.Sprintf(":%d\r\n", int(ttl.Seconds()))
		}
	}
	return ":-1\r\n"
}

func (db *Database) zadd(parts []string) string {
	// Implementation for ZADD command
	return "-ERR ZADD command is not implemented yet\r\n"
}

func (db *Database) zrange(parts []string) string {
	// Implementation for ZRANGE command
	return "-ERR ZRANGE command is not implemented yet\r\n"
}

func (db *Database) expired(key string) bool {
	expiry, ok := db.expiry[key]
	if !ok {
		return false
	}
	return time.Now().After(expiry)
}

func handleConnection(conn net.Conn, db *Database) {
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

func main() {
	db := NewDatabase()

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, db)
	}
}
