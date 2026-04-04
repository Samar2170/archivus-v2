package service

import (
	"archivus-v2/config"
	"archivus-v2/internal/utils"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const ThumbnailDir = ".thumbnails"

// prepareThumbnailDir creates the directory if needed and attempts to ensure it is writable.
// MkdirAll returns nil when the directory already exists regardless of its permissions,
// so we always attempt chmod to fix directories created by a previous run as a different user.
// chmod is best-effort: if the directory is owned by another user (e.g. root) we log a warning
// and continue — the subsequent write will surface the actual permission error if the dir is
// truly not writable.
func prepareThumbnailDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.Chmod(dir, 0755); err != nil {
		utils.LogError("prepareThumbnailDir",
			fmt.Sprintf("cannot set permissions on %s — if thumbnails fail, fix ownership with: sudo chown -R $USER %s", dir, dir),
			err)
	}
	return nil
}

func EnsureThumbnail(relPath string) (string, error) {
	switch {
	case IsImage(relPath):
		return ensureImageThumbnail(relPath)
	case IsVideo(relPath):
		return ensureVideoThumbnail(relPath)
	case IsPDF(relPath):
		return ensurePDFThumbnail(relPath)
	default:
		return "", nil
	}
}

func ensureImageThumbnail(relPath string) (string, error) {
	fullSourcePath := filepath.Join(config.Config.BaseDir, relPath)
	thumbRelPath := filepath.Join(ThumbnailDir, relPath)
	fullThumbPath := filepath.Join(config.Config.BaseDir, thumbRelPath)

	if _, err := os.Stat(fullThumbPath); err == nil {
		return thumbRelPath, nil
	}

	if err := prepareThumbnailDir(filepath.Dir(fullThumbPath)); err != nil {
		return "", utils.HandleError("ensureImageThumbnail", "Failed to prepare thumbnail directory", err)
	}

	srcImage, err := imaging.Open(fullSourcePath)
	if err != nil {
		utils.LogError("ensureImageThumbnail", "Failed to open image", err)
		return "", nil
	}

	thumbnail := imaging.Fit(srcImage, 200, 200, imaging.Lanczos)
	if err = imaging.Save(thumbnail, fullThumbPath); err != nil {
		utils.LogError("ensureImageThumbnail", "Failed to save thumbnail", err)
		return "", nil
	}

	return thumbRelPath, nil
}

func ensureVideoThumbnail(relPath string) (string, error) {
	fullSourcePath := filepath.Join(config.Config.BaseDir, relPath)
	thumbRelPath := filepath.Join(ThumbnailDir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".jpg")
	fullThumbPath := filepath.Join(config.Config.BaseDir, thumbRelPath)

	if _, err := os.Stat(fullThumbPath); err == nil {
		return thumbRelPath, nil
	}

	if err := prepareThumbnailDir(filepath.Dir(fullThumbPath)); err != nil {
		return "", utils.HandleError("ensureVideoThumbnail", "Failed to prepare thumbnail directory", err)
	}

	// Pipe stdout to avoid ffmpeg's image2 muxer failing on paths with spaces
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", fullSourcePath,
		"-ss", "00:00:01",
		"-vframes", "1",
		"-vf", "scale=200:-2",
		"-f", "mjpeg",
		"pipe:1",
	)
	jpegData, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			utils.LogError("ensureVideoThumbnail", "Failed to generate video thumbnail", fmt.Errorf("%w\nstderr: %s", err, exitErr.Stderr))
		} else {
			utils.LogError("ensureVideoThumbnail", "Failed to generate video thumbnail", err)
		}
		return "", nil
	}

	if err := os.WriteFile(fullThumbPath, jpegData, 0644); err != nil {
		utils.LogError("ensureVideoThumbnail", "Failed to write video thumbnail", err)
		return "", nil
	}

	return thumbRelPath, nil
}

func ensurePDFThumbnail(relPath string) (string, error) {
	fullSourcePath := filepath.Join(config.Config.BaseDir, relPath)
	thumbRelPath := filepath.Join(ThumbnailDir, strings.TrimSuffix(relPath, filepath.Ext(relPath))+".jpg")
	fullThumbPath := filepath.Join(config.Config.BaseDir, thumbRelPath)

	if _, err := os.Stat(fullThumbPath); err == nil {
		return thumbRelPath, nil
	}

	if err := prepareThumbnailDir(filepath.Dir(fullThumbPath)); err != nil {
		return "", utils.HandleError("ensurePDFThumbnail", "Failed to prepare thumbnail directory", err)
	}

	// Render first page to a temp JPEG, then resize with imaging
	tmpFile := fullThumbPath + ".tmp.jpg"
	defer os.Remove(tmpFile)

	cmd := exec.Command("gs",
		"-dNOPAUSE", "-dBATCH",
		"-sDEVICE=jpeg",
		"-r72",
		"-dFirstPage=1", "-dLastPage=1",
		"-sOutputFile="+tmpFile,
		fullSourcePath,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		utils.LogError("ensurePDFThumbnail", "Failed to render PDF", fmt.Errorf("%w\noutput: %s", err, out))
		return "", nil
	}

	img, err := imaging.Open(tmpFile)
	if err != nil {
		utils.LogError("ensurePDFThumbnail", "Failed to open PDF render", err)
		return "", nil
	}

	thumbnail := imaging.Fit(img, 200, 200, imaging.Lanczos)
	if err = imaging.Save(thumbnail, fullThumbPath); err != nil {
		utils.LogError("ensurePDFThumbnail", "Failed to save PDF thumbnail", err)
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

func IsVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v":
		return true
	}
	return false
}

func IsPDF(filename string) bool {
	return strings.ToLower(filepath.Ext(filename)) == ".pdf"
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
