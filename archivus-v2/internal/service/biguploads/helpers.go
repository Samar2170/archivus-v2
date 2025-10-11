package biguploads

import (
	"archivus-v2/config"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	WarkHeader   = "X-Wark"
	IdxHeader    = "X-Idx"
	ChHashHeader = "X-Chunk-Hash"
)

type UploadInitRequest struct {
	FileName    string   `json:"file_name"`
	FileSize    int64    `json:"file_size"`
	ChunkHashes []string `json:"chunk_hashes"`
	ModTime     int64    `json:"mod_time"`
}

type Session struct {
	Wark          string       `json:"wark"`
	FileName      string       `json:"file_name"`
	Size          int64        `json:"size"`
	ChunkHashes   []string     `json:"chunk_hashes"`
	ChunksWritten map[int]bool `json:"chunks_written"`
	CreatedAt     int64        `json:"created_at"`
}

var (
	sessions   = map[string]*Session{}
	sessionsMu sync.Mutex
	inmemLocks = map[string]*sync.Mutex{} // wark -> mutex for file
)

func EnsureBigUploadDirs() error {
	for _, d := range []string{config.GetSessionsDir(), config.GetTmpDir()} {
		if err := os.MkdirAll(d, 0755); err != nil {
			return err
		}
	}
	return nil
}

func saveSession(sess *Session) error {
	b, _ := json.MarshalIndent(sess, "", " ")
	fn := filepath.Join(config.GetSessionsDir(), sess.Wark+".json")
	return os.WriteFile(fn, b, 0644)
}

func loadSession(wark string) (*Session, error) {
	fn := filepath.Join(config.GetSessionsDir(), wark+".json")
	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func computeWark(size int64, chunkHashes []string) string {
	h := sha512.New()
	h.Write([]byte(config.Config.ServerSalt))
	h.Write([]byte(strconv.FormatInt(size, 10)))
	for _, ch := range chunkHashes {
		h.Write([]byte(ch))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func getLock(wark string) *sync.Mutex {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	m, ok := inmemLocks[wark]
	if !ok {
		m = &sync.Mutex{}
		inmemLocks[wark] = m
	}
	return m
}

func computeChunkSize(total int64, chunks int) int {
	if chunks <= 0 {
		return 1 << 20
	}
	size := (total + int64(chunks) - 1) / int64(chunks)
	return int(size)
}

func chooseChunkSize(total int64) int {
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
	return int(base)
}
