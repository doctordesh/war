package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/doctordesh/war"
	"github.com/doctordesh/war/colors"
)

func main() {

	var delay int
	var match, exclude string
	var verbose, boring bool

	flag.Usage = func() {
		fmt.Println("Usage: war [options] <command-to-run>")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.IntVar(&delay, "delay", 100, "Time in milliseconds before running command. Events within the delay will reset the delay")
	flag.StringVar(&match, "match", "*", "Match files, separate with comma")
	flag.StringVar(&exclude, "exclude", "", "Pattern to exclude files, separate with comma")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&boring, "boring", false, "Boring (no colors) output")

	flag.Parse()

	colors.SetColoring(!boring)

	args := flag.Args()
	if len(args) == 0 {
		colors.Red("missing <command> argument")
		flag.Usage()
		os.Exit(2)
	}

	if len(args) > 1 {
		colors.Yellow("ignores arguments '%v'", strings.Join(args[1:], "', '"))
	}

	path, err := os.Getwd()
	if err != nil {
		colors.Red("could not find current working directory: %v", err)
		os.Exit(2)
	}

	matches := splitAndTrim(match)
	excludes := splitAndTrim(exclude)

	w := war.New(path, matches, excludes, args[0], os.Environ(), delay)
	w.Verbose = verbose

	err = w.WatchAndRun()
	if err != nil {
		colors.Red("could not start war: %v", err)
		os.Exit(2)
	}
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
