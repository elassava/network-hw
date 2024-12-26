package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)
func main() {
	// Start the server on port 9000
	dstream, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer dstream.Close()

	fmt.Printf("Server is running at %s\n", dstream.Addr().String())

	var wg sync.WaitGroup

	for {
		// Accept new client connections
		conn, err := dstream.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			handleClient(c)
		}(conn)
	}
}

// Handle each client connection
func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading username:", err)
		return
	}
	username = username[:len(username)-1] // Trim newline character

	fmt.Printf("User connected: %s (%s)\n", username, conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		message := scanner.Text()
		// Print messages to the console
		fmt.Printf("[%s] -> %s\n", username, message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from client:", err)
	}
	fmt.Printf("User %s disconnected.\n", username)
}
