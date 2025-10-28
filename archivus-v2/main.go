package main

import (
	"archivus-v2/config"
	"archivus-v2/internal/syncer"
	"archivus-v2/pkg/logging"
	"archivus-v2/server"
	"archivus-v2/shell"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

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

		run()
	case newUserCmd.Happened():
		fmt.Println("Creating new master user...")
		shell.NewUser()
	case featureCmd.Happened():
		fmt.Println("Adding new feature...")

	case toggleUserSettingsCmd.Happened():
		shell.ToggleUserSettings()
	case syncCmd.Happened():
		fmt.Println("Syncing files...")
		runSync()

	case cleanupSyncDirCmd.Happened():
		fmt.Println("Cleaning up sync directory...")
		syncer.CleanupDirQueue()
	default:
		fmt.Println("No command provided. Use -h for help.")
	}
	fmt.Println("Mode:", *m)
}
func runSync() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errs := syncer.Sync(ctx)

	for _, err := range errs {
		fmt.Println(err)
	}

	if len(errs) == 0 {
		fmt.Println("Files synced successfully!")
	}

	fmt.Println("Syncing files...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}

func run() {
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		fmt.Println("Starting server...")
		defer wg.Done()
		server.RunServer(ctx)
	}()
	wg.Add(1)
	go startCronServer(ctx, &wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Shutting down...")
	cancel()

	wg.Wait()
	fmt.Println("Exited cleanly.")
}
