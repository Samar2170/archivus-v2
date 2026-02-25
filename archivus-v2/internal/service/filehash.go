package service

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/zeebo/blake3"
)

func HashFileBlake3(path string) (string, int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	hasher := blake3.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return "", 0, err
	}

	sum := hasher.Sum(nil)
	return hex.EncodeToString(sum), size, nil
}

// TODO: complete this function
func CreateFileHash()
