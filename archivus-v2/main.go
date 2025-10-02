package main

import (
	"archivus-v2/server"
	"archivus-v2/shell"
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("archivus-v2", "A simple file archiver")
	serverCmd := parser.NewCommand("server", "Run the archivus server")
	newUserCmd := parser.NewCommand("new-user", "Create a new master user")
	featureCmd := parser.NewCommand("feature", "Add a new feature")
	toggleUserSettingsCmd := parser.NewCommand("user-settings", "Toggle user settings")
	m := parser.String("m", "mode", &argparse.Options{
		Required: false,
		Help:     "Mode: 'archive' or 'extract'",
	})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	switch {
	case serverCmd.Happened():
		fmt.Println("Starting server...")
		server.RunServer()
	case newUserCmd.Happened():
		fmt.Println("Creating new master user...")
		shell.NewUser()
	case featureCmd.Happened():
		fmt.Println("Adding new feature...")
	case toggleUserSettingsCmd.Happened():
		shell.ToggleUserSettings()
	default:
		fmt.Println("No command provided. Use -h for help.")
	}
	fmt.Println("Mode:", *m)
}
