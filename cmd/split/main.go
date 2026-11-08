// Command split divides a bill among people and prints each share.
//
//	split 100 ana ben cem
//	split 100 ana=3 ben=1 cem=2
//
// A name without =N has a weight of 1.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/Costroid/costroid-demo/split"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "usage: split TOTAL NAME[=WEIGHT]...")
		return 2
	}
	total, err := split.ParseCents(args[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	people, weights, err := parsePeople(args[1:])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	shares, err := split.Weighted(total, people, weights)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	for _, share := range shares {
		fmt.Fprintf(stdout, "%s\t%s\n", share.Name, share.Amount)
	}
	return 0
}

// parsePeople reads NAME or NAME=WEIGHT for each person, in order.
func parsePeople(args []string) ([]string, map[string]int, error) {
	people := make([]string, 0, len(args))
	weights := make(map[string]int, len(args))
	for _, arg := range args {
		name, text, found := strings.Cut(arg, "=")
		weight := 1
		if found {
			parsed, err := strconv.Atoi(text)
			if err != nil {
				return nil, nil, fmt.Errorf("split: %s has a weight that is not a whole number", name)
			}
			weight = parsed
		}
		people = append(people, name)
		weights[name] = weight
	}
	return people, weights, nil
}
