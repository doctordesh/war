package war

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/doctordesh/check"
	"github.com/fsnotify/fsnotify"
)

func TestNotify(t *testing.T) {
	var err error
	var event fsnotify.Event
	var ok bool
	var act Action
	var file *os.File
	var n int

	basePath := t.TempDir()
	fooBar := filepath.Join(basePath, "foo-bar")
	newFile := filepath.Join(fooBar, "new-file")

	notify, err := fsnotify.NewWatcher()
	check.OK(t, err)

	t.Logf("Adds        %s to watcher", basePath)
	err = notify.Add(basePath)
	check.OK(t, err)

	t.Logf("Creates     %s", fooBar)
	err = os.Mkdir(fooBar, 0777)
	check.OK(t, err)

	event, ok = <-notify.Events
	check.Assert(t, ok)
	check.Equals(t, fsnotify.Create, event.Op)
	check.Equals(t, fooBar, event.Name)

	act, err = DecideAction(event, basePath, nil, nil, os.DirFS(basePath))
	check.OK(t, err)
	check.Equals(t, ActionAdd, act)

	t.Logf("Adds        %s to watcher", fooBar)
	err = notify.Add(fooBar)
	check.OK(t, err)

	t.Logf("Create at   %s", newFile)
	file, err = os.OpenFile(newFile, os.O_CREATE|os.O_RDWR, 0777)
	check.OK(t, err)

	n, err = file.Write([]byte("lorem"))
	check.OK(t, err)
	check.Equals(t, len("lorem"), n)

	t.Logf("File closed %s", newFile)
	err = file.Close()
	check.OK(t, err)

	t.Log("waiting for events")
	event, ok = <-notify.Events
	check.Assert(t, ok)
	check.Equals(t, fsnotify.Create, event.Op)
	check.Equals(t, newFile, event.Name)

	// Why do I need this? Is this to make the file system sync up?
	time.Sleep(time.Millisecond)

	file, err = os.OpenFile(newFile, os.O_RDWR, 0777)
	check.OK(t, err)
	n, err = file.Write([]byte("ipsum"))
	check.OK(t, err)
	check.Equals(t, len("ipsum"), n)

	err = file.Close()
	check.OK(t, err)

	t.Log("waiting for events")
	event, ok = <-notify.Events
	check.Assert(t, ok)
	check.Equals(t, fsnotify.Write, event.Op)
	check.Equals(t, newFile, event.Name)
}
