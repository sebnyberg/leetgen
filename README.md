# leetgen

Leetcode implementation and test case generator for Go (WIP)

## Todo

- [x] Create Leetcode client for fetching question contents
- [x] Parse function definition from code sample
- [x] Tokenize question contents, parse examples, combine with function definition
  to parse contents of inputs and outputs to their corresponding types.
- [x] Output stub based on question contents and function definition
- [x] Add generator flag
- [ ] Remove dependency on testify/require
- [ ] Add support for ListNode, TreeNode types
- [ ] Generate package from leetcode question
- [ ] Generate test cases based on data copied from submission results
- [ ] Generate tree of packages based on contest
