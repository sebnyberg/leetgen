package leetcode

import (
	"errors"
	"fmt"

	"golang.org/x/exp/slices"
)

type Problem struct {
	Fn ProblemFunc
}

func GetProblem(titleSlug string) (Problem, error) {
	descr, err := getProblemDescriptor(titleSlug)
	if err != nil {
		return Problem{}, err
	}

	// Find Go snippet, parse the function
	goIdx := slices.IndexFunc(descr.CodeSnippets,
		func(e codeSnippetDescriptor) bool {
			return e.LangSlug == "golang"
		})
	if goIdx == -1 {
		return Problem{}, errors.New("could not find go code snippet")
	}
	goSnippet := descr.CodeSnippets[goIdx]
	goFn, err := parseProblemFunc(goSnippet.Code)
	if err != nil {
		return Problem{}, fmt.Errorf("failed to parse function in sample code, %w", err)
	}

	fmt.Println(descr.ExampleTestcases)

	// Todo: parse test cases based on problem func

	p := Problem{
		Fn: goFn,
	}
	return p, nil
}
