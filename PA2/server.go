package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	//Listen
	dstream, err := net.Listen("tcp", "localhost:9000")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer dstream.Close()

	serverAddr := dstream.Addr().String()
	fmt.Printf("Launching TCP server at %s\n", serverAddr)

	//accept
	con, err := dstream.Accept() //con= connection
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()

	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(con)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()

		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	writer.Flush()

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
