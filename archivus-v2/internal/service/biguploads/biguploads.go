package biguploads

import (
	"archivus-v2/config"
	"archivus-v2/pkg/response"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func InitBigUpload(w http.ResponseWriter, r *http.Request) {
	var req UploadInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestResponse(w, "invalid request")
		return
	}
	if req.FileSize < 0 || len(req.FileName) == 0 || len(req.ChunkHashes) == 0 {
		response.BadRequestResponse(w, "invalid request")
		return
	}
	if len(req.ChunkHashes) > config.MaxChunks {
		response.BadRequestResponse(w, "too many chunks")
		return
	}

	wark := computeWark(req.FileSize, req.ChunkHashes)
	sess := &Session{
		Wark:          wark,
		FileName:      filepath.Base(req.FileName),
		Size:          req.FileSize,
		ChunkHashes:   req.ChunkHashes,
		ChunksWritten: map[int]bool{},
		CreatedAt:     time.Now().Unix(),
	}
	if exisiting, err := loadSession(wark); err == nil {
		sess = exisiting
	} else {
		tmpPath := filepath.Join(config.GetTmpDir(), sess.Wark+config.TmpSuffix)
		if _, err := os.Stat(tmpPath); errors.Is(err, os.ErrNotExist) {
			f, err := os.Create(tmpPath)
			if err != nil {
				response.InternalServerErrorResponse(w, "failed to create temp file")
				return
			}
			if err := f.Truncate(sess.Size); err != nil {
				f.Close()
				response.InternalServerErrorResponse(w, "failed to allocate space for temp file")
				return
			}
			f.Close()
		}
		if err := saveSession(sess); err != nil {
			response.InternalServerErrorResponse(w, "failed to save session")
			return
		}
	}

	needed := make([]int, 0, len(sess.ChunkHashes))
	for i := range sess.ChunkHashes {
		if !sess.ChunksWritten[i] {
			needed = append(needed, i)
		}
	}
	resp := map[string]interface{}{
		"wark":          sess.Wark,
		"needed_chunks": needed,
	}
	response.JSONResponse(w, resp)
}

func ChunkHandler(w http.ResponseWriter, r *http.Request) {
	wark := r.Header.Get(WarkHeader)
	idxS := r.Header.Get(IdxHeader)
	chHash := r.Header.Get(ChHashHeader)
	if wark == "" || idxS == "" || chHash == "" {
		response.BadRequestResponse(w, "missing required headers")
		return
	}
	idx, err := strconv.Atoi(idxS)
	if err != nil {
		response.BadRequestResponse(w, "invalid idx header")
		return
	}

	sessionsMu.Lock()
	sess, ok := sessions[wark]
	sessionsMu.Unlock()
	if !ok {
		if s, err := loadSession(wark); err == nil {
			sess = s
			sessionsMu.Lock()
			sessions[wark] = sess
			sessionsMu.Unlock()
		} else {
			response.BadRequestResponse(w, "invalid wark")
			return
		}
	}
	if idx < 0 || idx >= len(sess.ChunkHashes) {
		response.BadRequestResponse(w, "invalid idx")
		return
	}
	expected := strings.ToLower(sess.ChunkHashes[idx])
	if strings.ToLower(chHash) != expected {
		response.BadRequestResponse(w, "chunk hash does not match expected")
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		response.InternalServerErrorResponse(w, "failed to read body")
		return
	}
	sum := sha512.Sum512(bodyBytes)
	if hex.EncodeToString(sum[:]) != expected {
		response.BadRequestResponse(w, "chunk hash does not match body")
		return
	}

	tmpPath := filepath.Join(config.GetTmpDir(), sess.Wark+config.TmpSuffix)
	lock := getLock(wark)
	lock.Lock()
	defer lock.Unlock()

	f, err := os.OpenFile(tmpPath, os.O_WRONLY, 0644)
	if err != nil {
		response.InternalServerErrorResponse(w, "failed to open temp file")
		return
	}
	defer f.Close()

	chunkSize := computeChunkSize(sess.Size, len(sess.ChunkHashes))
	offset := int64(idx) * int64(chunkSize)
	if _, err := f.WriteAt(bodyBytes, offset); err != nil {
		response.InternalServerErrorResponse(w, "failed to write chunk to temp file")
		return
	}
	sess.ChunksWritten[idx] = true
	if err := saveSession(sess); err != nil {
		response.InternalServerErrorResponse(w, "failed to save session")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func FinaliseHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Wark string `json:"wark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequestResponse(w, "invalid request")
		return
	}
	if req.Wark == "" {
		response.BadRequestResponse(w, "missing wark")
		return
	}
	sess, err := loadSession(req.Wark)
	if err != nil {
		response.BadRequestResponse(w, "session not found")
		return
	}
	missing := []int{}
	for i := range sess.ChunkHashes {
		if !sess.ChunksWritten[i] {
			missing = append(missing, i)
		}
	}
	if len(missing) > 0 {
		response.BadRequestResponse(w, "not all chunks uploaded")
		return
	}

	tmpPath := filepath.Join(config.GetTmpDir(), sess.Wark+config.TmpSuffix)
	finalPath := filepath.Join(config.Config.BaseDir, sess.FileName)

	if _, err := os.Stat(finalPath); err == nil {
		finalPath = filepath.Join(config.Config.BaseDir, sess.Wark+"_"+sess.FileName)
	}

	if err := os.Rename(tmpPath, finalPath); err != nil {
		response.InternalServerErrorResponse(w, "failed to move file to final destination")
		return
	}

	os.Remove(filepath.Join(config.GetSessionsDir(), req.Wark+".json"))

	sessionsMu.Lock()
	delete(sessions, req.Wark)
	delete(inmemLocks, req.Wark)
	sessionsMu.Unlock()

	response.JSONResponse(w, map[string]string{"status": "ok", "path": finalPath})
}
