// hafize sanlı (220717050) & ela semra sava (220717026)

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Prompt the user for the filename
	fmt.Print("Connected to server...\nInput File Name: ")
	var filename string
	fmt.Scanln(&filename)

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Connect to the server
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Send file line by line
	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error sending data to server:", err)
			return
		}
		writer.Flush()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
}
