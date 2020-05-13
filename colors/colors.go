package colors

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

var au = aurora.NewAurora(true)

func SetColoring(b bool) {
	au = aurora.NewAurora(b)
}

func Yellow(str string, parts ...interface{}) {
	fmt.Printf("%s %s\n", au.BrightYellow("[WAR]"), fmt.Sprintf(str, parts...))
}

func Red(str string, parts ...interface{}) {
	fmt.Printf("%s %s\n", au.BrightRed("[WAR]"), fmt.Sprintf(str, parts...))
}

func Blue(str string, parts ...interface{}) {
	fmt.Printf("%s %s\n", au.BrightCyan("[WAR]"), fmt.Sprintf(str, parts...))
}

func Green(str string, parts ...interface{}) {
	fmt.Printf("%s %s\n", au.BrightGreen("[WAR]"), fmt.Sprintf(str, parts...))
}
