package main

import (
	"archivus-v2/config"
	"archivus-v2/internal/syncer"
	"archivus-v2/pkg/logging"
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

	syncCmd := parser.NewCommand("sync", "Sync files")
	cleanupSyncDirCmd := parser.NewCommand("cleanup-sync-queue", "Cleanup sync directory")
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	config.SetupConfig(*m)
	logging.SetupLogging()

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
	case syncCmd.Happened():
		fmt.Println("Syncing files...")
		errs := syncer.Sync()
		for _, err := range errs {
			fmt.Println(err)
		}
	case cleanupSyncDirCmd.Happened():
		fmt.Println("Cleaning up sync directory...")
		syncer.CleanupDirQueue()
	default:
		fmt.Println("No command provided. Use -h for help.")
	}
	fmt.Println("Mode:", *m)
}
