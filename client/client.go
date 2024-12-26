package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Send the username to the server
	_, err = conn.Write([]byte(username + "\n"))
	if err != nil {
		fmt.Println("Error sending username:", err)
		return
	}

	fmt.Println("Connected to the server. Type your messages below:")
	fmt.Println("To send DM: @username message")

	// Gelen mesajları dinlemek için goroutine
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Mesaj gönderme
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		if message == "quit" {
			return
		}

		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
