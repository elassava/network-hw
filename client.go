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
	var username string

	// Username validation loop
	for {
		fmt.Print("Enter your username: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		// Send the username to the server
		_, err = conn.Write([]byte(username + "\n"))
		if err != nil {
			fmt.Println("Error sending username:", err)
			return
		}

		// Wait for server response
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading server response:", err)
			return
		}
		response = strings.TrimSpace(response)

		if response == "USERNAME_TAKEN" {
			fmt.Println("Username already taken. Please choose another one.")
		} else if response == "USERNAME_ACCEPTED" {
			break
		}
	}

	fmt.Println("Connected to the server. Type your messages below:")
	fmt.Println("To send DM: [username] message")
	fmt.Println("To broadcast: [all] message")
	fmt.Println("To quit: quit\n")

	// goroutine to listen to incoming messages
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// Sending massage
	for {
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)

		//quit from server when users write quit
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
