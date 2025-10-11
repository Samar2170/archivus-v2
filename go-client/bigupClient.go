package main

import (
	"bufio"
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
	"strings"
	"sync"
	"time"
)

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

func uploadChunk(baseUrl, wark string, idx int, chunk []byte, hash string, headers map[string]string) error {
	url := baseUrl + UploadChunkUrl
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
		for k, v := range headers {
			req.Header.Set(k, v)
		}
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

func runBigUploadInteractive(baseUrl, filepath string) {
	username := getUserInput("Enter username: ", "")
	pin := getUserInput("Enter PIN: ", "")
	if len(username) < 3 || len(pin) != 6 {
		fmt.Println("Invalid username or PIN")
		return
	}
	lr := LoginRequest{
		Username: username,
		Pin:      pin,
	}
	err := DoBigUpload(baseUrl, filepath, lr)
	if err != nil {
		fmt.Println("Upload failed:", err)
	} else {
		fmt.Println("Upload succeeded")
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Pin      string `json:"pin"`
}

func login(client *http.Client, baseUrl string, lr LoginRequest) (string, error) {
	var err error
	req, err := http.NewRequest(
		"POST",
		baseUrl+LoginUrl,
		bytes.NewReader([]byte(fmt.Sprintf(`{"username":"%s","pin":"%s"}`, lr.Username, lr.Pin))),
	)
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	resp, err := MakeRequest(client, req, headers)
	if err != nil {
		return "", err
	}
	var loginResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		fmt.Println(loginResp)
		return "", fmt.Errorf("login failed")
	}
	fmt.Println("Login successful, token:", token)
	return token, nil
}

func initializeUpload(client *http.Client, filepath, baseUrl string, headers map[string]string) (initResponse, int64, int64, []string, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return initResponse{}, 0, 0, nil, err
	}
	totalSize := info.Size()
	chunkSize := chooseChunkSize(totalSize)
	fmt.Printf("Using chunk size: %d bytes\n", chunkSize)

	hashes, err := computeChunkHashes(filepath, chunkSize)
	if err != nil {
		return initResponse{}, 0, 0, nil, err
	}
	fmt.Printf("Computed %d chunk hashes\n", len(hashes))

	initReq := map[string]interface{}{
		"file_name":    info.Name(),
		"file_size":    totalSize,
		"chunk_hashes": hashes,
		"mod_time":     info.ModTime().Unix(),
	}
	body, _ := json.Marshal(initReq)
	req, err := http.NewRequest(
		"POST",
		baseUrl+InitBigUploadUrl,
		bytes.NewReader(body),
	)
	if err != nil {
		return initResponse{}, 0, 0, nil, err
	}
	resp, err := MakeRequest(client, req, headers)
	if err != nil {
		return initResponse{}, 0, 0, nil, err
	}
	var initResp initResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return initResponse{}, 0, 0, nil, err
	}
	resp.Body.Close()
	fmt.Println("wark:", initResp.Wark)
	return initResp, totalSize, chunkSize, hashes, nil

}

func DoBigUpload(baseUrl, filepath string, lr LoginRequest) error {
	client := &http.Client{}
	headers := make(map[string]string)
	token, err := login(client, baseUrl, lr)
	if err != nil {
		return err
	}
	headers["Authorization"] = "Bearer " + token

	initResp, totalSize, chunkSize, hashes, err := initializeUpload(client, filepath, baseUrl, headers)
	if err != nil {
		return err
	}
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
			if err := uploadChunk(baseUrl, initResp.Wark, idx, buf, hashes[idx], headers); err != nil {
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
	err = finalizeUpload(client, baseUrl, headers, initResp)
	return err
}

func finalizeUpload(client *http.Client, baseUrl string, headers map[string]string, initResp initResponse) error {
	finalReqBody := map[string]interface{}{
		"wark": initResp.Wark,
	}
	body, _ := json.Marshal(finalReqBody)
	finalReq, err := http.NewRequest(
		"POST",
		baseUrl+FinalizeUploadUrl,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	finalResp, err := MakeRequest(client, finalReq, headers)
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
