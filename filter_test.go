package war

import (
	"fmt"
	"testing"

	"github.com/doctordesh/check"
)

type matchcase struct {
	path string
	res  bool
}

// func TestFilterMatchAll(t *testing.T) {
// 	filter := Filter([]string{"*"})
// 	table := []matchcase{
// 		{"/home/thing", true},
// 		{"", true},
// 		{"/home/.git/index", true},
// 	}

// 	for k, row := range table {
// 		res := filter.Match(row.path)
// 		check.Assert(t, res == row.res, fmt.Sprintf("case %d failed", k))
// 	}
// }

func TestFilterMatchAll(t *testing.T) {
	filter := Filter{
		[]string{"/home/project", "/home/project/subdir"},
		[]string{"*"},
	}

	table := []matchcase{
		{"/home/project/file", false},
		{"/home/project", false},
		{"", false},
		{"/home/thing/index", false},
	}

	for k, row := range table {
		res := filter.Match(row.path)
		check.Assert(t, res == row.res, fmt.Sprintf("case %d failed", k))
	}
}
