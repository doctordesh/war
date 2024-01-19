package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/doctordesh/war"
	"github.com/doctordesh/war/colors"
)

var usage = func() {
	fmt.Println("Usage: war [options] <command-to-run>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	var err error
	var environment arrayArg
	var path, cwd string
	var boring bool
	var delay time.Duration

	flag.Var(&environment, "env", "Environment string with key=value pairs")
	flag.StringVar(&path, "path", ".", "Path to watch files from")
	flag.StringVar(&cwd, "cwd", ".", "Working directory of the command")
	flag.DurationVar(&delay, "delay", time.Millisecond*100, "Time before running command. Events within the delay will reset the delay")
	flag.BoolVar(&boring, "boring", false, "Boring (no colors) output")

	// usage of the program
	flag.Usage = usage

	flag.Parse()

	// set the coloring
	colors.SetColoring(!boring)

	args := flag.Args()
	if len(args) < 1 {
		colors.Red("missing <command> argument")
		flag.Usage()
		os.Exit(2)
	}

	binPath, err := exec.LookPath(args[0])
	if err != nil {
		colors.Red(err.Error())
		os.Exit(1)
	}

	rtpl := war.RunnableTemplate{
		BinPath: binPath,
		Args:    args,
		Env:     environment,
		Dir:     cwd,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}

	w := war.New(path, rtpl, time.Second)
	w.Verbose = true

	err = w.WatchAndRun()
	if err != nil {
		colors.Red("could not start war: %v", err)
		os.Exit(2)
	}
}

// arrayArg is a type to be able to pass multiple flags of the same name, and
// get them in a list. Only works with strings
type arrayArg []string

func (self *arrayArg) String() string {
	return strings.Join(*self, " ")
}

func (self *arrayArg) Set(value string) error {
	*self = append(*self, value)
	return nil
}
