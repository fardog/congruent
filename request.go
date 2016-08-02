package congruent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fardog/congruent/urljoin"
	"io/ioutil"
	"net/http"
	"strings"
)

// NewServer creates a new server definition
func NewServer(u string, h *http.Header) *Server {
	return &Server{BaseURI: u, Headers: h}
}

// Server represents a server that will be requested against
type Server struct {
	BaseURI string
	Headers *http.Header
}

// NewRequest creates a new request to be made against a Server
func NewRequest(m, p string, h *http.Header, b interface{}) *Request {
	return &Request{m, p, h, b}
}

// Request represents a test case to be run
type Request struct {
	Method  string
	Path    string
	Headers *http.Header
	Body    interface{}
}

// PrepareBody sets the body for a Request; can take a string, or any object which
// is treated as though it were JSON. TODO(nwittstock): support other things
func (r Request) PrepareBody() (*bytes.Buffer, error) {
	b := r.Body

	switch b.(type) {
	case string:
		return bytes.NewBuffer([]byte(b.(string))), nil
	default:
		body, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		return bytes.NewBuffer(body), nil
	}

}

// Do performs a Request and returns a Response
func (r Request) Do(s *Server) (*Response, error) {
	uri := urljoin.Join(s.BaseURI + r.Path)

	reqBody, err := r.PrepareBody()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, uri, reqBody)
	if err != nil {
		return nil, err
	}

	mergeHTTPHeaders(&req.Header, s.Headers, r.Headers)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{req, &resp.Header, body, resp.StatusCode}, nil
}

// Response represents a response from a server
type Response struct {
	Request    *http.Request
	Headers    *http.Header
	Body       []byte
	StatusCode int
}

type result struct {
	resp *Response
	err  error
}

// Servers is an array of Server pointers
type Servers []*Server

// Requests is an array of Request pointers
type Requests []*Request

// Request makes a Request against a list of servers, and returns responses
func (s Servers) Request(r *Request) (Responses, error) {
	var responses Responses
	var errors []string

	results := make(chan result, 4)

	for _, server := range s {
		go func(server *Server) {
			resp, err := r.Do(server)
			results <- result{resp, err}
		}(server)
	}

	for i := 0; i < len(s); i++ {
		result := <-results

		responses = append(responses, result.resp)

		if result.err != nil {
			errors = append(errors, result.err.Error())
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("got errors: %v", strings.Join(errors, "\t\n"))
	}

	return responses, nil
}
