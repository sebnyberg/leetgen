package leetcode

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/exp/slices"
)

type Problem struct {
	Fn       ProblemFunc
	Examples []Example
}

func GetProblem(titleSlug string) (*Problem, error) {
	descr, err := getProblemDescriptor(titleSlug)
	if err != nil {
		return nil, err
	}

	// Find Go snippet, parse the function
	goIdx := slices.IndexFunc(descr.CodeSnippets,
		func(e codeSnippetDescriptor) bool {
			return e.LangSlug == "golang"
		})
	if goIdx == -1 {
		return nil, errors.New("could not find go code snippet")
	}
	goSnippet := descr.CodeSnippets[goIdx]
	goFn, err := parseProblemFunc(goSnippet.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to parse function in sample code, %w", err)
	}

	// Parse examples from problem contents
	examples, err := parseExamplesFromContents(descr.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse problem examples, %w", err)
	}

	p := &Problem{
		Fn:       goFn,
		Examples: examples,
	}
	return p, nil
}

// WriteStub writes a problem stub to the provided writer.
func (p *Problem) WriteStub(w io.Writer, packageName string) (err error) {
	buf := make([]byte, 0, 128)

	// Indent helpers
	var nindent int
	indent := func() {
		if len(buf) > nindent {
			buf = buf[:nindent]
		}
		buf = append(buf, '\t')
		nindent++
	}
	outdent := func() {
		if nindent == 0 {
			panic("invalid indent")
		}
		nindent--
		buf = buf[:nindent]
	}

	// Write helpers
	write := func(bs []byte) {
		_, writeErr := w.Write(bs)
		if writeErr != nil && err == nil {
			err = writeErr
		}
	}
	writes := func(s string, args ...any) {
		buf = append(buf, fmt.Sprintf(s, args...)...)
		buf = append(buf, '\n')
		write(buf)
		buf = buf[:nindent]
	}
	writesin := func(s string, args ...any) {
		writes(s, args...)
		indent()
	}
	writesout := func(s string, args ...any) {
		outdent()
		writes(s, args...)
	}
	writeTestCase := func(inputs []string, output string) {
		vars := make([]string, len(inputs), len(inputs)+1)
		for i := range inputs {
			vars[i] = varToGoFormat(p.Fn.params[i].typ, inputs[i])
		}
		if p.Fn.retType != paramTypNone {
			vars = append(vars, varToGoFormat(p.Fn.retType, output))
		} else {
			vars = append(vars, varToGoFormat(p.Fn.params[0].typ, output))
		}
		single := fmt.Sprintf("{%v},", strings.Join(vars, ", "))
		if len(single) < 65 {
			writes(single)
			return
		}
		// Write on multiple rows
		writesin("{")
		for _, goInput := range vars {
			writes("%v,", goInput)
		}
		writesout("},")
	}

	writes("package %v\n", packageName)
	writesin("import (")
	writes(`"fmt"`)
	writes("\"testing\"\n")

	writes(`"github.com/stretchr/testify/require"`)
	writesout(")")

	writesin("func Test_%v(t *testing.T) {", p.Fn.name)
	writesin("type testCase struct {")
	for _, param := range p.Fn.params {
		writes("%v %v", param.name, param.typ.String())
	}
	if p.Fn.retType != paramTypNone {
		writes("want %v", p.Fn.retType.String())
	} else {
		// Guess that the ret type must correspond ot first input
		writes("want %v", p.Fn.params[0].typ.String())
	}
	writesout("}\n")
	writesin("testCases := []testCase{")
	for _, ex := range p.Examples {
		writeTestCase(ex.Inputs, ex.Output)
	}
	writesout("}\n")
	writesin("for i, tc := range testCases {")
	writesin(`t.Run(fmt.Sprintf("TestCase %%v", i), func(t *testing.T) {`)
	argNames := make([]string, len(p.Fn.params))
	for i := range argNames {
		argNames[i] = "tc." + p.Fn.params[i].name
	}
	funcCall := fmt.Sprintf("%v(%v)", p.Fn.name, strings.Join(argNames, ", "))
	if p.Fn.retType != paramTypNone {
		writes("require.Equal(t, tc.want, %v)", funcCall)
	} else {
		writes(funcCall)
		writes("require.Equal(t, tc.want, tc.%v)", p.Fn.params[0].name)
	}
	writesout("})")
	writesout("}")
	writesout("}\n")
	writes(p.Fn.Signature())
	return err
}

var bracketReplacer = strings.NewReplacer("[", "{", "]", "}")
var spacer = strings.NewReplacer(",", ", ")

func varToGoFormat(typ paramTyp, s string) string {
	// For any type of slice, prefix the "[]" or "[][]", the type, then replace
	// all "[" with "{" and vice versa
	s = spacer.Replace(s)
	if typ&paramTypSlice > 0 {
		return typ.String() + bracketReplacer.Replace(s)
	}
	if typ&paramTypSliceOfSlices > 0 {
		return typ.String() + bracketReplacer.Replace(s)
	}
	return s
}
