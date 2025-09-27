package main

import (
	"archivus/shell"
	"log"
	"os"
	"sync"
	"time"
)

var Commands = map[string][2]string{
	"setup":    {"setup", "Initialize the application with necessary configurations."},
	"sync":     {"sync", "Synchronize the application with the existing data."},
	"start":    {"start", "Start the application server."},
	"new-user": {"new-user", "Create a new user in the application."},
}

func Run(args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "setup":
			// This command is used to set up the application.
			log.Println("Setting up the application...")
		case "sync":
			// syncer.ManualSync()
		case "start":
			log.Println("Starting the application server...")
			// server.RunServer()
		case "new-master":
			shell.NewMasterUser()
		case "new-user":
			shell.Setup()
			shell.NewUser()
		}
	} else {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			// server.RunServer()
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			time.Sleep(2 * time.Second) // Wait for server to start
			startCronServer()
			wg.Done()
		}()
		wg.Wait()
	}
}

func main() {
	// This is the main entry point of the application.
	// It can be used to initialize and run the service.
	// Currently, it serves as a placeholder for future logic.
	args := os.Args[1:]
	Run(args)

}
