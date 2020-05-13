package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/doctordesh/war"
	"github.com/doctordesh/war/colors"
)

var usage = func() {
	fmt.Println("Usage: war [options] <command-to-run>")
	fmt.Println("Options:")
	flag.PrintDefaults()
}

func main() {
	// flag variables
	var delay int
	var match, exclude string
	var verbose, boring bool

	// setup flags
	flag.IntVar(&delay, "delay", 100, "Time in milliseconds before running command. Events within the delay will reset the delay")
	flag.StringVar(&match, "match", "*", "Match files, separate with comma")
	flag.StringVar(&exclude, "exclude", "", "Pattern to exclude files, separate with comma")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&boring, "boring", false, "Boring (no colors) output")

	// usage of the program
	flag.Usage = usage

	flag.Parse()

	// set the coloring
	colors.SetColoring(!boring)

	// Validate that we got the path
	args := flag.Args()
	if len(args) == 0 {
		colors.Red("missing <command> argument")
		flag.Usage()
		os.Exit(2)
	}

	// Warn about unnecessary arguments
	if len(args) > 1 {
		colors.Yellow("ignores arguments '%v'", strings.Join(args[1:], "', '"))
	}

	// Get the current working directory to know what to watch
	path, err := os.Getwd()
	if err != nil {
		colors.Red("could not find current working directory: %v", err)
		os.Exit(2)
	}

	// Extract matches and excludes from the string format
	matches := splitAndTrim(match)
	excludes := splitAndTrim(exclude)

	// Build program
	w := war.New(path, matches, excludes, args[0], os.Environ(), delay)
	w.Verbose = verbose

	// Setup signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Run
	go func() {
		err = w.WatchAndRun()
		if err != nil {
			colors.Red("could not start war: %v", err)
			os.Exit(2)
		}
	}()

	<-sigs

	fmt.Println()
	colors.Blue("keyboard interrupt detected, quiting...")
	os.Exit(0)
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	res := []string{}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])

		if parts[i] != "" {
			res = append(res, parts[i])
		}
	}

	return res
}
