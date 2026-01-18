package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/utils"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const ThumbnailDir = ".thumbnails"

func EnsureThumbnail(relPath string) (string, error) {
	if !IsImage(relPath) {
		return "", nil
	}

	fullSourcePath := filepath.Join(config.Config.BaseDir, relPath)
	thumbRelPath := filepath.Join(ThumbnailDir, relPath)
	fullThumbPath := filepath.Join(config.Config.BaseDir, thumbRelPath)

	// Check if thumbnail already exists
	if _, err := os.Stat(fullThumbPath); err == nil {
		return thumbRelPath, nil
	}

	// Ensure directory exists
	thumbDir := filepath.Dir(fullThumbPath)
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return "", utils.HandleError("EnsureThumbnail", "Failed to create thumbnail directory", err)
	}

	// Generate thumbnail
	srcImage, err := imaging.Open(fullSourcePath)
	if err != nil {
		// Log error but don't fail the request, just return empty string
		fmt.Printf("Failed to open image for thumbnail: %v\n", err)
		return "", nil
	}

	// Resize to 200x200 maintaining aspect ratio
	// imaging.Thumbnail handles cropping, Fit handles resizing within box
	thumbnail := imaging.Fit(srcImage, 200, 200, imaging.Lanczos)

	// Save
	err = imaging.Save(thumbnail, fullThumbPath)
	if err != nil {
		fmt.Printf("Failed to save thumbnail: %v\n", err)
		return "", nil
	}

	return thumbRelPath, nil
}

func IsImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	}
	return false
}

// Helper to get image dimensions if needed, though mostly for validation
func getImageConfig(path string) (image.Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return image.Config{}, err
	}
	defer file.Close()
	cfg, _, err := image.DecodeConfig(file)
	return cfg, err
}
