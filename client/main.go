package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	// 1. username
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	fmt.Fprint(conn, username)

	// 2. listen for messages
	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				return
			}
			fmt.Print(msg)
		}
	}()

	// 3. send messages
	for {
		text, _ := reader.ReadString('\n')
		fmt.Fprint(conn, text)
	}
}
