package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
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

	reader := bufio.NewReader(conn)
	for {
		// Read the first line to determine the type of command
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection", err)
			return
		}
		fmt.Printf("Received: %s", line)

		if strings.HasPrefix(line, "*") {
			// It's an array, parse the number of elements
			count, _ := strconv.Atoi(strings.TrimSpace(line[1:]))

			// Read each element
			var command, value string
			for i := 0; i < count; i++ {
				// Read the length line
				_, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading length:", err)
					return
				}

				// Read the actual data
				data, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading data:", err)
					return
				}

				if i == 0 {
					command = strings.TrimSpace(data)
				} else if i == 1 {
					value = strings.TrimSpace(data)
				}
			}

			// Process the command
			switch command {
			case "PING":
				conn.Write([]byte("+PONG\r\n"))
			case "ECHO":
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
			default:
				conn.Write([]byte("-ERR Unknown command\r\n"))
			}
		} else {
			fmt.Println("Unexpected format:", line)
			conn.Write([]byte("-ERR Protocol error\r\n"))
		}
	}
}
