package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/sebnyberg/leetgen/leetcode"
)

var problemName = flag.String("problem", "strange-printer", "problem stub, e.g. two-sum")

func main() {
	flag.Parse()
	if len(*problemName) == 0 {
		fmt.Println(`missing problem name, provide with e.g. -problem "two-sum"`)
		os.Exit(1)
	}
	problem, err := leetcode.GetProblem(*problemName)
	if err != nil {
		log.Fatalln(err)
	}
	problem.WriteStub(os.Stdout)
}
