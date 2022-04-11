package leetcode

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ProblemFunc struct {
	name    string
	params  []funcParam
	retType paramTyp
}

var funcRegex = regexp.MustCompile(`func (?P<Name>[^\s)]+)\s*\((?P<Params>[^)]+)\)\s*(?P<RetType>[^\n\s{]+)?`)

func parseProblemFunc(s string) (ProblemFunc, error) {
	parts := funcRegex.FindStringSubmatch(s)
	if len(parts) != 4 {
		return ProblemFunc{}, fmt.Errorf("parse function err, %w", errInvalidFormat)
	}
	var f ProblemFunc
	f.name = parts[funcRegex.SubexpIndex("Name")]

	// Parse params
	paramStrs := strings.Split(parts[funcRegex.SubexpIndex("Params")], ",")
	for i := range paramStrs {
		paramStrs[i] = strings.TrimSpace(paramStrs[i])
	}
	f.params = make([]funcParam, len(paramStrs))
	for i, paramStr := range paramStrs {
		var err error
		if f.params[i], err = parseParam(paramStr); err != nil {
			return ProblemFunc{}, fmt.Errorf("parse param failed, %w", err)
		}
	}

	// Parse return type (may be none)
	f.retType = parseParamTyp(parts[funcRegex.SubexpIndex("RetType")])

	return f, nil
}

func (f ProblemFunc) Output() string {
	paramStrs := make([]string, len(f.params))
	for i, p := range f.params {
		paramStrs[i] = fmt.Sprintf("%v %v", p.name, p.typ.String())
	}
	return fmt.Sprintf("func %v(%v) %v {\n\t\n}",
		f.name,
		strings.Join(paramStrs, ", "),
		f.retType.String(),
	)
}

type funcParam struct {
	name string
	typ  paramTyp
}

var errInvalidFormat = errors.New("invalid format")
var errUnknownType = errors.New("unknown type")

func parseParam(s string) (funcParam, error) {
	name, typeDef, found := strings.Cut(s, " ")
	if !found {
		return funcParam{}, fmt.Errorf(
			"%w: function param must contain two parts separated by space, was '%v'",
			errInvalidFormat, s,
		)
	}
	name = strings.TrimSpace(name)
	typeDef = strings.TrimSpace(typeDef)
	typ := parseParamTyp(typeDef)
	if typ == paramTypNone {
		return funcParam{}, errUnknownType
	}

	return funcParam{
		name: name,
		typ:  typ,
	}, nil
}
