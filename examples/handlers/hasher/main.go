package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Create a new scanner to read from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Loop indefinitely to process input continuously
	for {
		// Read the next line from stdin
		if scanner.Scan() {
			password := scanner.Text()

			// Generate bcrypt hash of the password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error hashing password:", err)
				continue
			}

			// Write the hashed password to stdout
			fmt.Println(string(hashedPassword))
		} else {
			// Check for errors during scanning
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Error reading input:", err)
			}
			// Break the loop if there's no more input (EOF)
			break
		}
	}
}

