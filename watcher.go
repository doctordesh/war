package war

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/doctordesh/war/colors"
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	directory string
	match     []string
	exclude   []string
	verbose   bool
}

type Watcher interface {
	Watch() (<-chan string, error)
	SetVerboseLogging(bool)
}

func NewWatcher(directory string, match []string, exclude []string) Watcher {
	return &watcher{directory, match, exclude, false}
}

func (w *watcher) SetVerboseLogging(b bool) {
	w.verbose = b
}

func (w *watcher) Watch() (<-chan string, error) {
	c := make(chan string)

	notify, err := fsnotify.NewWatcher()
	if err != nil {
		return c, err
	}

	// Find all sub directories
	dirs, err := w.allDirs(w.directory)
	if err != nil {
		return c, err
	}

	dirs = append(dirs, w.directory)

	if w.verbose {
		for _, d := range dirs {
			colors.Blue("watching %s", d)
		}
	} else {
		colors.Blue("watching %s", w.directory)
	}

	go func() {
		for {
			select {
			case event, ok := <-notify.Events:
				if !ok {
					colors.Red("unexpected error from notifier, could not read events")
					os.Exit(2)
				}

				// Not interested in chmods
				if event.Op&fsnotify.Chmod == fsnotify.Chmod {
					continue
				}

				// Not interested in removals
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					continue
				}

				// Ignore dot files and files inside dot-directories (.git)
				if w.isDotFile(event.Name) {
					continue
				}

				// Extract filename relative to root watching directory
				filename := strings.Replace(event.Name, w.directory+"/", "", -1)

				// Dir, add to watchers
				if w.isDir(event.Name) {

					// If it's not a create, it's irrelevant
					if event.Op&fsnotify.Create != fsnotify.Create {
						continue
					}

					// if it's a dot file, it's irrelevant
					if w.isDotFile(event.Name) {
						continue
					}

					// Add the new dir to the paths that we watch
					colors.Blue("new directory detected %s", filename)
					err = notify.Add(event.Name)
					if err != nil {
						colors.Red("could not add %s to notifier: %v", filename, err)
						os.Exit(2)
					}

					continue
				}

				if w.isMatch(event.Name, w.match) == false {
					continue
				}

				colors.Blue("file changed: %s", filename)

				c <- filename

			case err, _ := <-notify.Errors:
				colors.Red("unexpected error from notifier: %v", err)
				os.Exit(2)
			}
		}
	}()

	for _, d := range dirs {
		err = notify.Add(d)
		if err != nil {
			return c, fmt.Errorf("watcher could not add directory %s to notifier: %w", d, err)
		}
	}

	return c, nil
}

func (w *watcher) allDirs(dir string) ([]string, error) {
	dirs := []string{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return dirs, err
	}

	for _, f := range files {
		if f.IsDir() == false {
			continue
		}

		// build subdir path
		subDirPath := filepath.Join(dir, f.Name())

		// filter . files
		if w.isDotFile(subDirPath) {
			continue
		}

		// add valid subdir
		dirs = append(dirs, subDirPath)

		// Fetch all sub dirs from the found dir
		subDirs, err := w.allDirs(subDirPath)
		if err != nil {
			return dirs, err
		}

		dirs = append(dirs, subDirs...)
	}

	return dirs, nil
}

func (w *watcher) isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func (w *watcher) isMatch(path string, match []string) bool {
	return true
}

func (w *watcher) isDotFile(path string) bool {
	parts := strings.Split(path, "/")
	for i := range parts {
		if strings.HasPrefix(parts[i], ".") {
			return true
		}
		if parts[i] == "__pycache__" {
			return true
		}
	}

	return false
}
