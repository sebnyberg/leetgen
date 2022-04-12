package leetcode

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseContents(t *testing.T) {
	for i, tc := range []struct {
		problemFn string
		s         string
		want      []testCase
	}{
		{
			"func twoSum(nums []int, target int) []int {\n\n}",
			"<p>Given an array of integers <code>nums</code>&nbsp;and an integer <code>target</code>, return <em>indices of the two numbers such that they add up to <code>target</code></em>.</p>\n\n<p>You may assume that each input would have <strong><em>exactly</em> one solution</strong>, and you may not use the <em>same</em> element twice.</p>\n\n<p>You can return the answer in any order.</p>\n\n<p>&nbsp;</p>\n<p><strong>Example 1:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [2,7,11,15], target = 9\n<strong>Output:</strong> [0,1]\n<strong>Explanation:</strong> Because nums[0] + nums[1] == 9, we return [0, 1].\n</pre>\n\n<p><strong>Example 2:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [3,2,4], target = 6\n<strong>Output:</strong> [1,2]\n</pre>\n\n<p><strong>Example 3:</strong></p>\n\n<pre>\n<strong>Input:</strong> nums = [3,3], target = 6\n<strong>Output:</strong> [0,1]\n</pre>\n\n<p>&nbsp;</p>\n<p><strong>Constraints:</strong></p>\n\n<ul>\n\t<li><code>2 &lt;= nums.length &lt;= 10<sup>4</sup></code></li>\n\t<li><code>-10<sup>9</sup> &lt;= nums[i] &lt;= 10<sup>9</sup></code></li>\n\t<li><code>-10<sup>9</sup> &lt;= target &lt;= 10<sup>9</sup></code></li>\n\t<li><strong>Only one valid answer exists.</strong></li>\n</ul>\n\n<p>&nbsp;</p>\n<strong>Follow-up:&nbsp;</strong>Can you come up with an algorithm that is less than&nbsp;<code>O(n<sup>2</sup>)&nbsp;</code>time complexity?",
			[]testCase{
				{
					input:      []any{[]int{2, 7, 11, 15}, 9},
					wantOutput: any([]int{0, 1}),
				},
			},
		},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			// fn, err := parseproblemfunc(tc.problemfn)
			// require.noerror(t, err)
			// p := problem{
			// 	fn: fn,
			// }

			res, err := parseExamplesFromContents(tc.s)
			require.NoError(t, err)
			if len(tc.want) != len(res) {
				log.Fatalf("[%v]: result len (%v) != expected len (%v)", i, len(res), len(tc.want))
			}
			for j := range tc.want {
				require.ElementsMatch(t, tc.want[j], res[j], fmt.Sprint(j))
			}
		})
	}
}
