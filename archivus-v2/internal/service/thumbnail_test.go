package service

import (
	"archivus-v2/config"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureThumbnail(t *testing.T) {
	// Setup temporary directory for BaseDir
	tmpDir, err := os.MkdirTemp("", "archivus_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Mock Config
	config.Config = &config.Configuration{
		BaseDir: tmpDir,
	}

	// Create a dummy image structure
	// Create "uploads" folder or similar if needed, but we can just use root of tmpDir as source
	// Test case: image in root
	imageName := "test_image.png"
	fullImagePath := filepath.Join(tmpDir, imageName)

	createTestImage(t, fullImagePath)

	// Test EnsureThumbnail
	thumbRelPath, err := EnsureThumbnail(imageName)
	if err != nil {
		t.Fatalf("EnsureThumbnail failed: %v", err)
	}

	if thumbRelPath == "" {
		t.Errorf("Expected thumbnail path, got empty string")
	}

	expectedRelPath := filepath.Join(ThumbnailDir, imageName)
	if thumbRelPath != expectedRelPath {
		t.Errorf("Expected path %s, got %s", expectedRelPath, thumbRelPath)
	}

	// Verify file exists
	fullThumbPath := filepath.Join(tmpDir, thumbRelPath)
	if _, err := os.Stat(fullThumbPath); os.IsNotExist(err) {
		t.Errorf("Thumbnail file does not exist at %s", fullThumbPath)
	}

	// Test case 2: Non-image file
	textName := "test.txt"
	os.WriteFile(filepath.Join(tmpDir, textName), []byte("hello"), 0644)
	thumbRelPath, err = EnsureThumbnail(textName)
	if err != nil {
		t.Errorf("EnsureThumbnail failed for text file: %v", err)
	}
	if thumbRelPath != "" {
		t.Errorf("Expected empty string for text file, got %s", thumbRelPath)
	}

	// Test case 3: Nested directory
	subDir := filepath.Join(tmpDir, "sub", "folder")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	subImageName := filepath.Join("sub", "folder", "image.jpg")
	createTestImage(t, filepath.Join(tmpDir, subImageName))

	thumbRelPath, err = EnsureThumbnail(subImageName)
	if err != nil {
		t.Fatalf("EnsureThumbnail nested failed: %v", err)
	}

	expectedSubRelPath := filepath.Join(ThumbnailDir, subImageName)
	if thumbRelPath != expectedSubRelPath {
		t.Errorf("Expected nested path %s, got %s", expectedSubRelPath, thumbRelPath)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, thumbRelPath)); os.IsNotExist(err) {
		t.Errorf("Nested thumbnail file not found")
	}
}

func createTestImage(t *testing.T, path string) {
	width := 400
	height := 400
	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.White)
		}
	}

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create image file: %v", err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode image: %v", err)
	}
}
