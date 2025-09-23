package war

import "fmt"

type Action int

const (
	ActionIgnore = Action(0)
	ActionRun    = Action(1)
	ActionAdd    = Action(2)
)

// String ...
func (a Action) String() string {
	switch a {
	case ActionIgnore:
		return "Ignore"
	case ActionRun:
		return "Run"
	case ActionAdd:
		return "Add"
	}

	panic(fmt.Sprintf("unknwon action value '%d'", a))
}
