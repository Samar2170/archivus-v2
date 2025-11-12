package main

import (
	"archivus-v2/internal/syncer"
	"archivus-v2/pkg/logging"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

func getCronServer(ctx context.Context) *gocron.Scheduler {
	t := time.Now()
	logging.AuditLogger.Info().Msgf("Starting cron server at %s", t.Format(time.RFC3339))
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Hour().Do(func() {
		// err := image.MarkImages()
		// if err != nil {
		// 	logging.Errorlogger.Error().Msgf("Error in Marking images: %s", err.Error())
		// }
	})
	s.Every(1).Hour().Do(func() {
		// err := image.CompressImages(config.CompressionQuality)
		// if err != nil {
		// 	logging.Errorlogger.Error().Msgf("Error in Compressing images: %s", err.Error())
		// }
	})
	s.Every(1).Hour().Do(func() {
		fmt.Println("Syncing files...")
		errs := syncer.Sync(ctx)
		for _, err := range errs {
			logging.Errorlogger.Error().Msgf("Error in syncing files: %s", err.Error())
		}
	})

	return s
}

func startCronServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	s := getCronServer(ctx)
	s.StartAsync()
	fmt.Printf("Cron server started at %s\n", time.Now().Format(time.RFC3339))

	<-ctx.Done()
	s.Stop()
	log.Println("Cron server stopped")
}
