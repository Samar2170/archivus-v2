package client

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	MaxParallel = 6
	MaxRetries  = 3
)

type initResponse struct {
	Wark   string `json:"wark"`
	Needed []int  `json:"needed"`
}

func chooseChunkSize(total int64) int64 {
	base := int64(1 << 20)
	maxChunk := int64(32 << 20)
	chunks := (total + base - 1) / base
	for chunks > 256 && base < maxChunk {
		base *= 2
		chunks = (total + base - 1) / base
	}
	if base > maxChunk {
		base = maxChunk
	}
	return base
}

func computeChunkHashes(path string, chunkSize int64) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	info, _ := f.Stat()
	total := info.Size()
	num := int((total + chunkSize - 1) / chunkSize)
	hashes := make([]string, num)
	buf := make([]byte, chunkSize)
	for i := 0; i < num; i++ {
		need := int64(chunkSize)
		if int64(i)*chunkSize+need > total {
			need = total - int64(i)*chunkSize
		}
		n, err := io.ReadFull(f, buf[:need])
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			return nil, err
		}
		sum := sha512.Sum512(buf[:n])
		hashes[i] = hex.EncodeToString(sum[:])
	}
	return hashes, nil
}

func uploadChunk(wark string, idx int, chunk []byte, hash string) error {
	url := ""
	for attempt := 1; attempt <= MaxRetries; attempt++ {
		req, err := http.NewRequest(
			"POST",
			url,
			bytes.NewReader(chunk),
		)
		if err != nil {
			return err
		}
		req.Header.Set("X-Wark", wark)
		req.Header.Set("X-Idx", strconv.Itoa(idx))
		req.Header.Set("X-Chunk-Hash", hash)
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			if attempt == MaxRetries {
				return err
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		b, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			if attempt == MaxRetries {
				return errors.New("upload failed: " + string(b))
			}
			time.Sleep(time.Duration(attempt) * time.Second)
			continue
		}
		return nil
	}
	return nil
}

func BigUploadTest(filepath string) error {
	info, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	totalSize := info.Size()
	chunkSize := chooseChunkSize(totalSize)
	fmt.Printf("Using chunk size: %d bytes\n", chunkSize)

	hashes, err := computeChunkHashes(filepath, chunkSize)
	if err != nil {
		return err
	}
	fmt.Printf("Computed %d chunk hashes\n", len(hashes))

	initReq := map[string]interface{}{
		"filename":     info.Name(),
		"filesize":     totalSize,
		"chunk_hashes": hashes,
		"mod_time":     info.ModTime().Unix(),
	}
	body, _ := json.Marshal(initReq)
	resp, err := http.Post(
		"",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	var initResp initResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return err
	}
	resp.Body.Close()
	fmt.Println("wark:", initResp.Wark)
	needednMap := make(map[int]bool)
	for _, n := range initResp.Needed {
		needednMap[n] = true
	}
	f, _ := os.Open(filepath)
	defer f.Close()
	sem := make(chan struct{}, MaxParallel)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var firstErr error

	for idx := 0; idx < len(hashes); idx++ {
		if !needednMap[idx] {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()
			offset := int64(idx) * chunkSize
			toRead := chunkSize
			if offset+toRead > totalSize {
				toRead = totalSize - offset
			}
			buf := make([]byte, toRead)
			_, err := f.ReadAt(buf, offset)
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			if err := uploadChunk(initResp.Wark, idx, buf, hashes[idx]); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				return
			}
			fmt.Printf("uploaded chunk %d\n", idx)
		}(idx)
	}
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}

	finalReq := map[string]interface{}{
		"wark": initResp.Wark,
	}
	body, _ = json.Marshal(finalReq)
	finalResp, err := http.Post(
		"",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	var fResp map[string]interface{}
	if err := json.NewDecoder(finalResp.Body).Decode(&fResp); err != nil {
		return err
	}
	finalResp.Body.Close()
	fmt.Println("Upload complete:", fResp)
	return nil
}
