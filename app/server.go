package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		// go handleConnection(conn)    ####### for handling multiple connections simultaneously
		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from connection", err)
		return
	}
	fmt.Printf("input is: %s", input)
	command_lenght, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Printf("input is: %s", command_lenght)

	command, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Printf("input is: %s", command)
	if strings.TrimSpace(command) == "PING" {
		// fmt.Fprint(conn, "+PONG\r\n")
		conn.Write([]byte("+PONG\r\n"))
	}
}
