package main

import (
	"archivus/config"
	"archivus/internal/helpers"
	"archivus/internal/service/image"
	"archivus/internal/syncer"
	"archivus/pkg/logging"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron"
)

func getCronServer() *gocron.Scheduler {
	t := time.Now()
	logging.AuditLogger.Info().Msgf("Starting cron server at %s", t.Format(time.RFC3339))
	s := gocron.NewScheduler(time.UTC)
	s.Every(2).Hour().Do(func() {
		err := image.MarkImages()
		if err != nil {
			logging.Errorlogger.Error().Msgf("Error in Marking images: %s", err.Error())
		}
	})
	s.Every(1).Hour().Do(func() {
		err := image.CompressImages(config.CompressionQuality)
		if err != nil {
			logging.Errorlogger.Error().Msgf("Error in Compressing images: %s", err.Error())
		}
	})
	s.Every(1).Hour().Do(func() {
		helpers.UpdateDirsData()
		helpers.UpdateUserDirsData()
		syncer.Sync()
	})

	return s
}

func startCronServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	s := getCronServer()
	go func() {
		s.StartBlocking()
	}()
	fmt.Printf("Cron server started at %s\n", time.Now().Format(time.RFC3339))
	<-stop

	fmt.Printf("Shutting down cron server at %s\n", time.Now().Format(time.RFC3339))
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Stop()
	log.Println("Cron server stopped")

}
