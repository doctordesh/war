package war

import "path/filepath"

type Filter struct {
	Directories []string
	Filters     []string
}

func (f Filter) Match(candidate string) bool {
	match := false
	for _, filter := range f {
		res, err := filepath.Match(filter, candidate)
		if err != nil {
			panic(err)
		}

		if res {
			match = true
			break
		}
	}

	return match
}
