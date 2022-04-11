package leetcode

import (
	"log"
	"strconv"
	"strings"
)

type paramTyp uint8

const (
	paramTypNone          = 0
	paramTypString        = 1 << 0
	paramTypBool          = 1 << 1
	paramTypInt           = 1 << 2
	paramTypInt64         = 1 << 3
	paramTypSlice         = 1 << 4
	paramTypSliceOfSlices = 1 << 5
)

func parseParamTyp(s string) paramTyp {
	s = strings.TrimSpace(s)

	var typ paramTyp
	if len(s) > 4 && s[:4] == "[][]" {
		typ |= paramTypSliceOfSlices
		s = s[4:]
	} else if len(s) > 2 && s[:2] == "[]" {
		typ |= paramTypSlice
		s = s[2:]
	}
	switch s {
	case "string":
		typ |= paramTypString
	case "int":
		typ |= paramTypInt
	case "int64":
		typ |= paramTypInt64
	case "bool":
		typ |= paramTypBool
	}
	return typ
}

func (t paramTyp) String() string {
	res := make([]byte, 0, 10)
	if t&paramTypSliceOfSlices > 0 {
		res = append(res, "[][]"...)
	} else if t&paramTypSlice > 0 {
		res = append(res, "[]"...)
	}
	switch {
	case t&paramTypBool > 0:
		res = append(res, "bool"...)
	case t&paramTypInt > 0:
		res = append(res, "int"...)
	case t&paramTypInt64 > 0:
		res = append(res, "int64"...)
	case t&paramTypString > 0:
		res = append(res, "string"...)
	}
	return string(res)
}

func (t paramTyp) parse(s string) (any, error) {
	s = strings.TrimSpace(s)

	if t&paramTypSliceOfSlices > 0 {
		// Parse each element as a slice
		t &^= paramTypSliceOfSlices
		t |= paramTypSlice
		els := strings.Split(s[1:len(s)-1], ",")
		res := make([]any, len(els))
		var err error
		for i, el := range els {
			res[i], err = t.parse(el)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}
	if t&paramTypSlice > 0 {
		els := strings.Split(s, ",")
		res := make([]any, len(els))
		t &^= paramTypSlice
		var err error
		for i, el := range els {
			res[i], err = t.parse(el)
			if err != nil {
				return nil, err
			}
		}
		return res, nil
	}
	switch t {
	case paramTypInt:
		x, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return x, nil
	case paramTypInt64:
		x, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return int64(x), nil
	case paramTypString:
		return s[1 : len(s)-1], nil
	case paramTypBool:
		panic("bool not implemented")
	default:
		log.Fatalf("unknown type %v", t)
		return nil, errUnknownType
	}
}
