package leetcode_test

import (
	"testing"

	"github.com/sebnyberg/leetgen/leetcode"
)

func TestClient(t *testing.T) {
	titleSlug := "two-sum"
	p, err := leetcode.GetProblem(titleSlug)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	_ = p
}
