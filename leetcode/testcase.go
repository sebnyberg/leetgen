package leetcode

import "fmt"

type testCase struct {
	input      []any
	wantOutput any
}

func newTestCase(fn ProblemFunc, inputs []string, output string) (testCase, error) {
	if len(inputs) != len(fn.params) {
		return testCase{}, fmt.Errorf(
			"invalid number of inputs (%v), expected %v",
			len(inputs), len(fn.params),
		)
	}
	args := make([]any, len(inputs))
	for i := range inputs {
		var err error
		args[i], err = fn.params[i].typ.parse(inputs[i])
		if err != nil {
			return testCase{}, err
		}
	}
	want, err := fn.retType.parse(output)
	if err != nil {
		return testCase{}, err
	}
	return testCase{
		input:      args,
		wantOutput: want,
	}, nil
}
