package leetcode

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Example struct {
	InputRaw    string
	Inputs      []string
	Output      string
	Explanation string
}

type exampleParser struct {
	z           *html.Tokenizer
	examples    []Example
	exampleRows []string
	rowBuf      []byte
	err         error
}

const (
	exampleInputIdx       = 0
	exampleOutputIdx      = 1
	exampleExplanationIdx = 2
)

type stateFn func(*exampleParser) stateFn

func parseExamplesFromContents(contents string) ([]Example, error) {
	z := html.NewTokenizerFragment(strings.NewReader(contents), "p")
	parser := &exampleParser{
		z: z,
	}
	// parse examples from problem contents
	for state := parseExamples; state != nil; {
		state = state(parser)
	}

	return parser.examples, parser.err
}

func parseExamples(p *exampleParser) stateFn {
	// Parse everything up to <pre>
	for {
		tt := p.z.Next()
		switch tt {
		case html.ErrorToken:
			if p.z.Err() != io.EOF {
				p.err = p.z.Err()
			}
			return nil
		case html.StartTagToken:
			if p.z.Token().String() == "<pre>" {
				return parseExample(p)
			}
		}
	}
}

func parseExample(p *exampleParser) stateFn {
	// Parse up to and including <strong>
	p.exampleRows = p.exampleRows[:0]
	for {
		tt := p.z.Next()
		switch tt {
		case html.ErrorToken:
			if p.z.Err() != io.EOF {
				p.err = p.z.Err()
			}
			return nil
		case html.StartTagToken:
			if p.z.Token().String() == "<strong>" {
				return parseExampleRowHeader(p)
			}
		case html.TextToken: // do nothing
		default:
			panic(p.z.Token().String())
		}
	}
}

func parseExampleRowHeader(p *exampleParser) stateFn {
	// Parse up to and including </strong>
	for {
		tt := p.z.Next()
		switch tt {
		case html.StartTagToken:
		case html.TextToken:
		case html.EndTagToken:
			s := p.z.Token().String()
			if s == "</strong>" {
				return parseExampleRowContents(p)
			}
		default:
			panic(p.z.Token().String())
		}
	}
}

func parseExampleRowContents(p *exampleParser) stateFn {
	// Parse up to and including <strong> or </pre>
	p.rowBuf = p.rowBuf[:0]
	for {
		tt := p.z.Next()
		switch tt {
		case html.StartTagToken:
			s := p.z.Token().String()
			p.exampleRows = append(p.exampleRows, string(p.rowBuf))
			if s == "<strong>" {
				return parseExampleRowHeader(p)
			}
		case html.EndTagToken:
			s := p.z.Token().String()
			if s == "</pre>" {
				p.exampleRows = append(p.exampleRows, string(p.rowBuf))
				return emitExample(p)
			}
			panic(s)
		case html.TextToken:
			s := strings.TrimSpace(p.z.Token().String())
			if s == "" {
				break
			}
			p.rowBuf = append(p.rowBuf, s...)
		default:
			panic(p.z.Token().String())
		}
	}
}

func emitExample(p *exampleParser) stateFn {
	if len(p.exampleRows) < 2 {
		p.err = errors.New("example did not contain at least two rows")
		return nil
	}
	if len(p.exampleRows) > 3 {
		p.err = errors.New("example contained more than three rows")
		return nil
	}

	var e Example

	// Convert example inputs to follow the same format as testCases
	parts := strings.Split(p.exampleRows[exampleInputIdx], ", ")
	e.Inputs = make([]string, len(parts))
	for i, part := range parts {
		_, right, ok := strings.Cut(part, " = ")
		if !ok {
			p.err = fmt.Errorf("invalid input variable format '%v'", p)
			return nil
		}
		e.Inputs[i] = strings.TrimSpace(right)
	}
	e.InputRaw = p.exampleRows[exampleInputIdx]
	e.Output = p.exampleRows[exampleOutputIdx]
	if len(p.exampleRows) == 3 {
		e.Explanation = p.exampleRows[exampleExplanationIdx]
	}

	p.examples = append(p.examples, e)
	return parseExamples
}
