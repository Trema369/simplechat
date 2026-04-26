package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mu      sync.Mutex
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	fmt.Println("Server started on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 1. First message = username
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	mu.Lock()
	clients[conn] = username
	mu.Unlock()

	broadcast(fmt.Sprintf("%s joined the chat\n", username), conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		msg = strings.TrimSpace(msg)
		fullMsg := fmt.Sprintf("%s: %s\n", username, msg)

		broadcast(fullMsg, conn)
	}

	mu.Lock()
	delete(clients, conn)
	mu.Unlock()

	broadcast(fmt.Sprintf("%s left the chat\n", username), conn)
}

func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message))
		}
	}
}
