package leetcode

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseType(t *testing.T) {
	for i, tc := range []struct {
		s       string
		want    paramTyp
		wantErr error
	}{
		{"int", paramTypInt, nil},
		{"[]int", paramTypInt | paramTypSlice, nil},
		{"[][]int", paramTypInt | paramTypSliceOfSlices, nil},
		{"string", paramTypString, nil},
		{"[]string", paramTypString | paramTypSlice, nil},
		{"[][]string", paramTypString | paramTypSliceOfSlices, nil},
		{"int64", paramTypInt64, nil},
		{"bool", paramTypBool, nil},
	} {
		t.Run(fmt.Sprintf("[%v] %v", i, tc.s), func(t *testing.T) {
			res, err := parseParam(tc.s)
			if tc.wantErr != nil {
				require.ErrorIs(t, err, tc.wantErr)
			} else {
				require.Equal(t, tc.want, res)
			}
		})
	}
}
