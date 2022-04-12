package leetcode

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProblem(t *testing.T) {
	fn, err := parseProblemFunc("func twoSum(a string, nums []int, target int) []int {\n\n}")
	require.NoError(t, err)
	p := Problem{
		Fn: fn,
		Examples: []Example{
			{
				InputRaw:    "",
				Inputs:      []string{`"abc"`, "[1,2,3,4,5]", "1"},
				Output:      "[1,2,3]",
				Explanation: "",
			},
		},
	}
	require.NoError(t, p.WriteStub(os.Stdout, "leetcode"))
}
