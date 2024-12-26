package main

import (
    "bufio"
    "fmt"
    "os"
    "strconv"
)

func main() {
    // Prompt the user for input and output filenames
    var inputFilename, outputFilename string
    fmt.Print("Enter the input filename: ")
    fmt.Scanln(&inputFilename)
    fmt.Print("Enter the output filename: ")
    fmt.Scanln(&outputFilename)

    // Open the input file for reading
    inputFile, err := os.Open(inputFilename)
    if err != nil {
        fmt.Println("Error opening input file:", err)
        return
    }
    defer inputFile.Close()

    // Create or open the output file for writing
    outputFile, err := os.Create(outputFilename)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer outputFile.Close()

    // Create buffered readers and writers
    scanner := bufio.NewScanner(inputFile)
    writer := bufio.NewWriter(outputFile)

    // Process each line in the input file
    for scanner.Scan() {
        line := scanner.Text()
        charCount := len(line) // Calculate the character count of the line

        // Prepare the output line with the character count prepended
        outputLine := strconv.Itoa(charCount) + " " + line + "\n"
        
        // Write the modified line to the output file
        _, err := writer.WriteString(outputLine)
        if err != nil {
            fmt.Println("Error writing to output file:", err)
            return
        }
    }

    // Check for scanner errors and flush the writer buffer
    if err := scanner.Err(); err != nil {
        fmt.Println("Error reading input file:", err)
    }
    writer.Flush()
}
