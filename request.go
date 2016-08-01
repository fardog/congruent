package congruent

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Headers represents HTTP header key/value pairs
type Headers map[string]string

// NewServer creates a new server definition
func NewServer(u string, h *Headers) *Server {
	return &Server{BaseURI: u, Headers: h}
}

// Server represents a server that will be requested against
type Server struct {
	BaseURI string
	Headers *Headers
}

// NewRequest creates a new request to be made against a Server
func NewRequest(m, p string, h *Headers, b interface{}) *Request {
	return &Request{m, p, h, b}
}

// Request represents a test case to be run
type Request struct {
	Method  string
	Path    string
	Headers *Headers
	Body    interface{}
}

// PrepareBody sets the body for a Request; can take a string, or any object which
// is treated as though it were JSON. TODO(nwittstock): support other things
func (r Request) PrepareBody() ([]byte, error) {
	b := r.Body

	switch b.(type) {
	case string:
		return []byte(b.(string)), nil
	default:
		body, err := json.Marshal(b)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

}

// Do performs a Request and returns a Response
func (r Request) Do(s *Server) (*Response, error) {
	uri := s.BaseURI + r.Path // TODO(nwittstock): proper path concat

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, uri, nil)
	if err != nil {
		return nil, err
	}

	if s.Headers != nil {
		for h, v := range *s.Headers {
			req.Header.Add(h, v)
		}
	}

	if r.Headers != nil {
		for h, v := range *r.Headers {
			req.Header.Add(h, v)
		}
	}

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
