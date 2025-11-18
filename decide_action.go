package war

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var ErrIsRelative = errors.New("path cannot be relative")
var ErrDifferentBasePaths = errors.New("path must be relative to base")

// var excludedParts = []string{
// 	".git",
// 	"__pycache__",
// 	".venv",
// }

// var excludedPaths = []string{
// 	"bin",
// }

func DecideAction(event fsnotify.Event, basePath string, excludedPaths, excludedPathParts []string, disk fs.FS) (Action, error) {
	if event.Op&fsnotify.Chmod == fsnotify.Chmod {
		return ActionIgnore, nil
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		return ActionIgnore, nil
	}

	path := filepath.Clean(event.Name)
	if filepath.IsAbs(path) == false {
		return ActionIgnore, ErrIsRelative
	}

	if strings.HasPrefix(path, basePath) == false {
		return ActionIgnore, ErrDifferentBasePaths
	}

	relPath, err := filepath.Rel(basePath, path)
	if err != nil {
		return ActionIgnore, fmt.Errorf("could not find relative path between '%s' and '%s': %w", basePath, event.Name, err)
	}

	if matchOneOf(relPath, excludedPathParts) {
		return ActionIgnore, nil
	}

	if matchSubPath(relPath, excludedPaths) {
		return ActionIgnore, nil
	}

	if isEmacsTempFile(relPath) {
		return ActionIgnore, nil
	}

	if isDir(relPath, disk) {
		if event.Op != fsnotify.Create {
			return ActionIgnore, nil
		}

		return ActionAdd, nil
	}

	return ActionRun, nil
}

// matchOneOf takes a path and a list pf pathParts. If one if the pathParts
// exist in the path it will return true, otherwise false.
func matchOneOf(path string, pathParts []string) bool {
	d, f := filepath.Split(path)

	if d == "" || d == "/" {
		return false
	}

	for _, name := range pathParts {
		if f == name {
			return true
		}
	}

	return matchOneOf(filepath.Clean(d), pathParts)
}

// matchSubPath returns true if a value in paths is a subpath to path
func matchSubPath(path string, paths []string) bool {
	for _, p := range paths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}

	return false
}

// isDir return true if the path is considered a directory according to the fs.FS
func isDir(path string, disk fs.FS) bool {
	file, err := disk.Open(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func isEmacsTempFile(path string) bool {
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".#") {
		return true
	}

	return false
}
