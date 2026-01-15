package main

import (
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

const (
	MaxParallel = 6
	MaxRetries  = 3

	InitBigUploadUrl  = "/bigupload/initiate/"
	UploadChunkUrl    = "/bigupload/chunk/"
	FinalizeUploadUrl = "/bigupload/finalize/"
	LoginUrl          = "/login"
)

type initResponse struct {
	Wark   string `json:"wark"`
	Needed []int  `json:"needed_chunks"`
}

func main() {
	parser := argparse.NewParser("archivus bigup Client", "A client for big file uploads")
	filepath := parser.String("f", "file", &argparse.Options{
		Required: true,
		Help:     "Path to the file to upload",
	})
	baseUrl := parser.String("u", "url", &argparse.Options{
		Required: false,
		Help:     "Base URL of the server (e.g., http://localhost:8080)",
		Default:  "http://localhost:8000",
	})
	debugMode := parser.Flag("d", "debug", &argparse.Options{
		Required: false,
		Help:     "Enable debug mode",
	})
	// usage
	// go-client -f /path/to/largefile -u http://localhost:8000
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}
	t0 := time.Now()
	runBigUploadInteractive(*baseUrl, *filepath, *debugMode)
	t1 := time.Now()
	fmt.Printf("Total upload time: %v\n", t1.Sub(t0))
}
