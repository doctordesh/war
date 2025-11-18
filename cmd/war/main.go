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

const VERSION = "0.1.42-rc14"

var usage = func() {
	fmt.Println("Usage: war [options] <command-to-run>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	var err error
	var environment, exclude arrayArg
	var cwd string
	var boring, version bool
	var delay, ignoreChangesFor time.Duration

	flag.Var(&environment, "env", "Environment string with key=value pairs")
	flag.Var(&exclude, "exclude", "Exclude changes on path, relative to the base path")
	flag.DurationVar(&delay, "delay", 0, "Time before running command")
	flag.DurationVar(&ignoreChangesFor, "ignore-changes-for", time.Millisecond*100, "Events within the specified time will be ignored and reset the delay")
	flag.BoolVar(&boring, "boring", false, "Boring (no colors) output")
	flag.BoolVar(&version, "version", false, "Print version and exit")

	// usage of the program
	flag.Usage = usage

	flag.Parse()

	if version {
		fmt.Printf("Version: %s\n", VERSION)
		os.Exit(0)
	}

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

	cwd, err = os.Getwd()
	if err != nil {
		colors.Red(err.Error())
		os.Exit(1)
	}

	rtpl := war.RunnableTemplate{
		BinPath:  binPath,
		Args:     args,
		Env:      environment,
		Excludes: exclude,
		Dir:      cwd,
		Stdout:   os.Stdout,
		Stderr:   os.Stderr,
	}

	w := war.New(cwd, rtpl, delay, ignoreChangesFor)
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
