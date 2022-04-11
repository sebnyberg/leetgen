package main

import (
	"fmt"
	"log"

	"github.com/sebnyberg/leetgen/leetcode"
)

func main() {
	// Todo: use CWD to fetch package name for output file
	// cwd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	problem, err := leetcode.GetProblem("search-insert-position")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(problem.Fn.Output())
}
