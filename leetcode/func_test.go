package leetcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseProblemFunc(t *testing.T) {
	for i, tc := range []struct {
		s       string
		want    ProblemFunc
		wantErr error
	}{
		{
			"func aFunc  (s string, ss [][]string,  b bool)   []int64 {",
			ProblemFunc{
				name: "aFunc",
				params: []funcParam{
					{"s", paramTypString},
					{"ss", paramTypString | paramTypSliceOfSlices},
					{"b", paramTypBool},
				},
				retType: paramTypInt64 | paramTypSlice,
			},
			nil,
		},
		{
			"func a(s string, b bool) int",
			ProblemFunc{
				name: "a",
				params: []funcParam{
					{"s", paramTypString},
					{"b", paramTypBool},
				},
				retType: paramTypInt,
			},
			nil,
		},
		{
			"func a(s string)",
			ProblemFunc{
				name: "a",
				params: []funcParam{
					{"s", paramTypString},
				},
				retType: paramTypNone,
			},
			nil,
		},
	} {
		t.Run(fmt.Sprintf("[%v] %v", i, tc.s), func(t *testing.T) {
			res, err := parseProblemFunc(tc.s)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.Equal(t, tc.want, res)
			}
		})
	}
}
