package leetcode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://leetcode.com"
const gqlURL = "https://leetcode.com/graphql"
const defaultTimeout = time.Minute * 1
const csrfCookieKey = "csrftoken"
const maxReadSize = 4 * 1024 * 1024 // 4MB

type client struct {
	csrftoken string
	cookie    *http.Cookie
}

func newClient() (client, error) {
	var c client
	var err error
	c.csrftoken, c.cookie, err = getCSRFToken()
	if err != nil {
		return c, fmt.Errorf("get csrf token failed, %w", err)
	}
	return c, nil
}

type codeSnippetDescriptor struct {
	Lang     string `json:"lang"`
	LangSlug string `json:"langSlug"`
	Code     string `json:"code"`
}
type problemDescriptor struct {
	QuestionID         string                  `json:"questionId"`
	QuestionFrontendID string                  `json:"questionFrontendId"`
	Title              string                  `json:"title"`
	TitleSlug          string                  `json:"titleSlug"`
	Content            string                  `json:"content"`
	ExampleTestcases   string                  `json:"exampleTestcases"`
	CodeSnippets       []codeSnippetDescriptor `json:"codeSnippets"`
}

func getProblemDescriptor(titleSlug string) (problemDescriptor, error) {
	c, err := newClient()
	if err != nil {
		return problemDescriptor{}, err
	}

	type respData struct {
		Question problemDescriptor `json:"question"`
	}
	type jsonResp struct {
		Data respData `json:"data"`
	}
	type jsonRequestVars struct {
		TitleSlug string `json:"titleSlug"`
	}
	type jsonRequest struct {
		OperationName string          `json:"operationName"`
		Variables     jsonRequestVars `json:"variables"`
		Query         string          `json:"query"`
	}

	const query = `
query questionData($titleSlug: String!) {
  question(titleSlug: $titleSlug) {
    exampleTestcases
    codeSnippets {
      lang
      langSlug
      code
    }
  }
}
`
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	// Setup request
	jsonReq := jsonRequest{
		OperationName: "questionData",
		Variables: jsonRequestVars{
			TitleSlug: titleSlug,
		},
		Query: query,
	}
	jsonBytes, err := json.Marshal(jsonReq)
	if err != nil {
		return problemDescriptor{}, fmt.Errorf("request json marshal err, %w", err)
	}

	// Make request
	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, gqlURL, bytes.NewReader(jsonBytes))
	if err != nil {
		return problemDescriptor{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", "https://leetcode.com/problems/two-sum/")
	req.Header.Set("x-csrftoken", c.csrftoken)
	req.AddCookie(c.cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return problemDescriptor{}, fmt.Errorf("gql request failed, %w", err)
	}

	// Read result
	var contents jsonResp
	lr := io.LimitReader(resp.Body, maxReadSize)
	if err := json.NewDecoder(lr).Decode(&contents); err != nil {
		return problemDescriptor{}, fmt.Errorf("decode Leetcode problem body err, %w", err)
	}

	return contents.Data.Question, nil
}

// getCSRFToken does a GET to Leetcode's website to retrieve the CSRF token
// through the cookie response. The cookie, and the the CSRF token must be
// provided in subsequent requests. CSRF is provided via the `x-csrftoken`
// header.
func getCSRFToken() (string, *http.Cookie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return "", nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", nil, err
	}
	if n := len(resp.Cookies()); n != 1 {
		return "", nil, fmt.Errorf("unexpected number of cookies: %v", n)
	}
	c := resp.Cookies()[0]
	if c.Name != csrfCookieKey {
		return "", nil, fmt.Errorf("unexpected cookie key: %v", c.Name)
	}
	return c.Value, c, nil
}
