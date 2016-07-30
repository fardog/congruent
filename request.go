package congruent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Headers represents HTTP header key/value pairs
type Headers map[string]string

// AuthDef allows passing auth data which will be encoded later; this
// can be provided in Headers, but provided here we'll do the encoding for you
type AuthDef struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a AuthDef) String() string {
	return a.BasicAuth()
}

// BasicAuth creates an authentication string suitable for use in a header
func (a AuthDef) BasicAuth() string {
	s := fmt.Sprintf("%s:%s", a.Username, a.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return fmt.Sprintf("Basic %s", encoded)
}

func NewServer(u string, h *Headers) *Server {
	return &Server{BaseURI: u, Headers: h}
}

// Server represents a server that will be requested against
type Server struct {
	BaseURI        string
	Headers        *Headers
	Authentication *AuthDef
}

func (s *Server) SetAuth(a *AuthDef) {
	s.Authentication = a
}

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

type Servers []*Server
type Requests []*Request

func (s Servers) Request(r *Request) (Responses, error) {
	var responses Responses

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
	}

	return responses, nil
}
