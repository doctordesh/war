package war

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

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
			log.Printf("watcher - watching %s", d)
		}
	}

	go func() {
		for {
			select {
			case event, ok := <-notify.Events:
				if !ok {
					log.Fatal("watcher - could not read events")
				}

				if w.verbose {
					log.Printf("watcher - event: %s", event)
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
					log.Printf("watcher - watching %s", event.Name)
					err = notify.Add(event.Name)
					if err != nil {
						log.Fatalf("watcher - could not add %s to notifier: %v", event.Name, err)
					}

					continue
				}

				if w.isMatch(event.Name, w.match) == false {
					continue
				}

				filename := strings.Replace(event.Name, w.directory+"/", "", -1)
				c <- filename

			case err, ok := <-notify.Errors:
				if !ok {
					log.Fatal("watcher - could not read errors")
				}

				log.Printf("watcher - error: %v", err)
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
