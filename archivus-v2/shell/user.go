package shell

import (
	"archivus-v2/internal"
	"archivus-v2/internal/auth"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getUserInput(prompt, defaultValue string) string {
	fmt.Print(prompt)

	// Check if stdin is a terminal
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error accessing stdin:", err)
		return defaultValue
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// Not a terminal (e.g., stdin redirected), return default
		fmt.Println("[No terminal detected, using default value]")
		return defaultValue
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return defaultValue
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}

	return input
}
func createNewUser(username, password, pin, email string, isMaster, createDir bool) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if len(pin) != 6 {
		return fmt.Errorf("pin must be exactly 6 digits long")
	}
	var ak string
	var err error
	if !isMaster {
		err = auth.CheckMasterUser()
		if err != nil {
			return fmt.Errorf("master user check failed: %w", err)
		}
	}
	ak, _, err = auth.CreateUser(username, password, pin, email, isMaster, createDir)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("User created successfully!")
	fmt.Println("Save this API key securely:")
	fmt.Println(ak)
	return nil
}

func NewUser() {
	internal.Setup(false)
	username := getUserInput("Enter username (at least 3 characters): ", "test")
	password := getUserInput("Enter password (at least 6 characters): ", "123456")
	pin := getUserInput("Enter pin (exactly 6 digits): ", "123456")
	email := getUserInput("Enter email address: ", "")
	isMasterUser := getUserInput("Is this a master user? (y/n): ", "y")
	var isMaster bool
	if strings.ToLower(isMasterUser) == "y" {
		isMaster = true
	} else {
		isMaster = false
	}
	fmt.Println("Creating new user...")
	fmt.Println("Params:")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", strings.Repeat("*", len(password)))
	fmt.Printf("Pin: %s\n", strings.Repeat("*", len(pin)))
	fmt.Printf("Email: %s\n", email)
	fmt.Println(username, password, pin, email, isMaster)
	var err error
	err = auth.SudoCheck()
	if err != nil {
		log.Fatalf("Sudo check failed: %v", err)
	}
	err = createNewUser(username, password, pin, email, isMaster, true)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
}
