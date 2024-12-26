// Client Code
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Set the username
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

	// Start sending messages
	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
