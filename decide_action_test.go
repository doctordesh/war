package war

import (
	"fmt"
	"io/fs"
	"testing"
	"time"

	"github.com/doctordesh/check"
	"github.com/fsnotify/fsnotify"
)

var disk mockFS

func TestIgnoreCmhods(t *testing.T) {
	var err error

	event := fsnotify.Event{Name: "not important", Op: fsnotify.Chmod}
	act, err := DecideAction(event, "", nil, nil, disk)

	check.OK(t, err)
	check.Assert(t, act == ActionIgnore)
}

func TestIgnoreRemovals(t *testing.T) {
	var err error

	event := fsnotify.Event{Name: "not important", Op: fsnotify.Remove}
	act, err := DecideAction(event, "", nil, nil, disk)

	check.OK(t, err)
	check.Assert(t, act == ActionIgnore)
}

func TestEventPathMustBeAbsolute(t *testing.T) {
	_, err := DecideAction(
		fsnotify.Event{Name: "local/path", Op: fsnotify.Create},
		"",
		nil, nil,
		disk,
	)

	check.Equals(t, ErrIsRelative, err)
}

func TestMustBePartOfBasePath(t *testing.T) {
	var err error
	_, err = DecideAction(
		fsnotify.Event{Name: "/absolute/but/wrong/base", Op: fsnotify.Create},
		"/absolute/with/right",
		nil, nil,
		disk,
	)

	check.Equals(t, ErrDifferentBasePaths, err)

	_, err = DecideAction(
		fsnotify.Event{Name: "/absolute/with/right/base", Op: fsnotify.Create},
		"/absolute/with/right",
		nil, nil,
		disk,
	)

	check.OK(t, err)
}

func TestIgnoreBasedOnParts(t *testing.T) {
	type row struct {
		ChangeAtPath string
		ShouldIgnore bool
	}

	base := "/foo/bar/.git/"
	table := []row{
		{base + "lorem/ipsum/.git", true},
		{base + "lorem/ipsum/git", false},
		{base + "./lorem/ipsum/.git", true},
		{base + "./lorem/ipsum/git", false},
		{base + "/lorem/ipsum/.git", true},
		{base + "/lorem/ipsum/git", false},
		{base + "/lorem/ipsum/.git/index", true},
		{base + "/lorem/ipsum/git/index", false},
		{base + "local/without/git/index", false},
		{base + "/local/without/git/index", false},
	}

	for _, row := range table {
		// res := matchOneOf(row.Path, []string{".git"})
		act, err := DecideAction(
			fsnotify.Event{Name: row.ChangeAtPath, Op: fsnotify.Create},
			base,
			nil,
			[]string{".git"},
			disk,
		)
		check.OKWithMessage(t, err, "for path %s", row.ChangeAtPath)
		check.AssertWithMessage(t, (act == ActionIgnore) == row.ShouldIgnore, "for path %s", row.ChangeAtPath)
	}
}

func TestIgnoresBasedOnPaths(t *testing.T) {
	type row struct {
		ChangeAtPath string
		ExcludePaths []string
		ShouldIgnore bool
	}

	base := "/proj/"
	table := []row{
		{base + "bin/executable", []string{"bin"}, true},
		{base + "sub/bin/executable", []string{"bin"}, false},
	}

	for _, row := range table {
		act, err := DecideAction(
			fsnotify.Event{Name: row.ChangeAtPath, Op: fsnotify.Create},
			base,
			row.ExcludePaths,
			nil,
			disk,
		)

		check.OKWithMessage(t, err, "for path %s", row.ChangeAtPath)
		check.AssertWithMessage(t, (act == ActionIgnore) == row.ShouldIgnore, "for path %s", row.ChangeAtPath)
	}
}

func TestIsDir(t *testing.T) {
	var disk fs.FS

	disk = &mockFS{
		files: map[string]mockFile{
			"foo/bar": {
				info: mockFileInfo{
					filename: "foo/bar",
					isDir:    true,
				},
			},
		},
	}

	type row struct {
		Path      string
		Op        fsnotify.Op
		ShouldAdd bool
	}

	table := []row{
		{"/foo/bar", fsnotify.Op(0), false},
		{"/foo/bar", fsnotify.Create, true},
	}

	for _, row := range table {
		res, err := DecideAction(fsnotify.Event{Name: row.Path, Op: row.Op}, "/", nil, nil, disk)
		check.OK(t, err)
		check.AssertWithMessage(t, (ActionAdd == res) == row.ShouldAdd, "for path %s and op %v", row.Path, row.Op)
	}
}

// ==================================================
//
// Mock Filesystem
//
// The goal is to be able to answer the 'IsDir' question
//
// ==================================================

type mockFS struct {
	files map[string]mockFile
}

func (m mockFS) Open(path string) (fs.File, error) {
	file, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("file '%s' not found", path)
	}

	return file, nil
}

type mockFile struct {
	info mockFileInfo
}

func (m mockFile) Stat() (fs.FileInfo, error) { return m.info, nil }
func (m mockFile) Read([]byte) (int, error)   { return 0, nil }
func (m mockFile) Close() error               { return nil }

type mockFileInfo struct {
	filename string
	isDir    bool
}

func (m mockFileInfo) Name() string       { return m.filename }
func (m mockFileInfo) Size() int64        { return 0 }
func (m mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m mockFileInfo) ModTime() time.Time { return time.Now() }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() any           { return nil }
