package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		// Accept an incoming connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection
		go handleConnection(conn) // for handling multiple connections simultaneously
		// handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	type storageItem struct {
		value      string
		expiryTime time.Time
	}

	storage := make(map[string]storageItem)

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection", err)
			return
		}
		fmt.Printf("Received: %s", line)

		if strings.HasPrefix(line, "*") {
			count, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
			var parts []string
			for i := 0; i < count; i++ {
				_, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading length:", err)
					return
				}
				data, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading data:", err)
					return
				}
				parts = append(parts, strings.TrimSpace(data))
			}

			command := strings.ToUpper(parts[0])
			switch command {
			case "PING":
				conn.Write([]byte("+PONG\r\n"))
			case "ECHO":
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(parts[1]), parts[1])))
			case "SET":
				key, value := parts[1], parts[2]
				item := storageItem{value: value}
				if len(parts) > 3 && strings.ToUpper(parts[3]) == "PX" {
					duration, _ := strconv.Atoi(parts[4])
					item.expiryTime = time.Now().Add(time.Duration(duration) * time.Millisecond)
				}
				storage[key] = item
				conn.Write([]byte("+OK\r\n"))
			case "GET":
				key := parts[1]
				if item, exists := storage[key]; exists {
					if item.expiryTime.IsZero() || time.Now().Before(item.expiryTime) {
						conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(item.value), item.value)))
					} else {
						delete(storage, key)
						conn.Write([]byte("$-1\r\n"))
					}
				} else {
					conn.Write([]byte("$-1\r\n"))
				}
			default:
				conn.Write([]byte("-ERR Unknown command\r\n"))
			}
		} else {
			fmt.Println("Unexpected format:", line)
			conn.Write([]byte("-ERR Protocol error\r\n"))
		}
	}
}
