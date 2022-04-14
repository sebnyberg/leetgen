package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/sebnyberg/leetgen/leetcode"
)

var problemName string
var packageName string

func main() {
	flag.StringVar(&packageName, "pkg", "", "package name, if unset defaults to e.g. 'p0123allinonestub'")
	flag.StringVar(&problemName, "p", "", "problem stub or URL, e.g. two-sum or https://leetcode.com/problems/spiral-matrix-ii/")
	flag.Parse()

	problemName, err := parseProblemName(problemName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	problem, err := leetcode.GetProblem(problemName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	problem.WriteStub(os.Stdout, "abc123")
}

func parseProblemName(problemName string) (string, error) {
	if len(problemName) == 0 {
		return "", errors.New(`missing problem name, provide with e.g. -problem "two-sum"`)
	}

	// If problemName is a URL, keep only the path
	if u, err := url.Parse(problemName); err != nil {
		problemName = u.Path
	}
	problemName = strings.TrimSuffix(problemName, "/")
	// Assume last part of the path is the problem name
	parts := strings.Split(problemName, "/")
	problemName = parts[len(parts)-1]
	return problemName, nil
}
