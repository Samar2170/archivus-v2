package synkr

import (
	"strings"
)

var skipDirsList = map[string]bool{
	"venv":         true,
	".git":         true,
	"__pycache__":  true,
	"snap":         true,
	"node_modules": true,
	"bin":          true,
	"build":        true,
}

var skipFilesList = map[string]bool{
	".DS_Store":   true,
	".gitkeep":    true,
	".gitignore":  true,
	"Thumbs.db":   true,
	"desktop.ini": true,
	"favicon.ico": true,
	"changelog":   true,
}

func shouldScanDir(dir string) bool {
	if dir == "" || dir[0] == '.' || dir[0] == '_' {
		return false
	}
	dirNameSplit := strings.Split(dir, "/")
	lastElement := dirNameSplit[len(dirNameSplit)-1]
	if lastElement[0] == '.' {
		return false
	}
	if _, ok := skipDirsList[lastElement]; ok {
		return false
	}
	return true
}

func shouldScanFile(file string) bool {
	if file == "" || file[0] == '.' || file[0] == '_' {
		return false
	}
	fileNameSplit := strings.Split(file, "/")
	lastElement := fileNameSplit[len(fileNameSplit)-1]
	if lastElement[0] == '.' {
		return false
	}
	if _, ok := skipFilesList[lastElement]; ok {
		return false
	}

	return true
}

func formatErrors(errs []error) string {
	var errStr string
	for _, err := range errs {
		errStr += err.Error() + "\n"
	}
	return errStr
}
