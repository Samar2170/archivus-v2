package shell

import (
	"archivus/internal/auth"
	"archivus/internal/service"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

func Setup() {
	service.Setup(false)
}

func getUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	value, _ := reader.ReadString('\n')
	value = strings.ReplaceAll(value, "\n", "")
	return value
}

func CreateNewUser(username, password, pin, email string, master bool) error {
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}
	if len(pin) != 6 {
		return fmt.Errorf("pin must be exactly 6 digits long")
	}
	ak, _, err := auth.CreateUser(username, password, pin, email, master)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("User created successfully!")
	fmt.Println("Save this API key securely:")
	fmt.Println(ak)
	return nil

}

func sudoCheck() error {
	cmd := exec.Command("sudo", "-v")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func NewMasterUser() {
	u, err := user.Current()
	if err != nil {
		log.Fatalf("Error fetching current user: %v", err)
	}
	err = sudoCheck()
	if err != nil {
		log.Fatalf("This operation requires sudo privileges: %v", err)
	}
	fmt.Println("Current OS User:", u.Username)
	username := getUserInput("Enter master username (at least 3 characters): ")
	password := getUserInput("Enter master password (at least 6 characters): ")
	pin := getUserInput("Enter master pin (exactly 6 digits): ")
	email := getUserInput("Enter master email address: ")
	fmt.Println("Creating new master user...")
	fmt.Println("Params:")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", strings.Repeat("*", len(password)))
	fmt.Printf("Pin: %s\n", strings.Repeat("*", len(pin)))
	fmt.Printf("Email: %s\n", email)

	err = CreateNewUser(username, password, pin, email, true)
	if err != nil {
		log.Fatalf("Error creating master user: %v", err)
	}
}

func NewUser() {
	username := getUserInput("Enter username (at least 3 characters): ")
	password := getUserInput("Enter password (at least 6 characters): ")
	pin := getUserInput("Enter pin (exactly 6 digits): ")
	email := getUserInput("Enter email address: ")
	fmt.Println("Creating new user...")
	fmt.Println("Params:")
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", strings.Repeat("*", len(password)))
	fmt.Printf("Pin: %s\n", strings.Repeat("*", len(pin)))
	fmt.Printf("Email: %s\n", email)

	err := CreateNewUser(username, password, pin, email, false)
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
}
