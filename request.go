package congruent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// NewRequest creates a request given a RequestDef and a ServerDef
func NewRequest(sd ServerDef, rd RequestDef) *Request {
	return &Request{
		rd.Method,
		sd.BaseURI + rd.Path, // TODO(nwittstock): join as a URL properly
		sd.Headers,
		nil,
	}
}

// Request represents a request to be made against a server
type Request struct {
	Method  string
	URI     string
	Headers HeadersDef
	body    []byte
}

// Response represents a response from a server
type Response struct {
	Headers    HeadersDef
	Body       []byte
	StatusCode int
}

// SetBody sets the body for a Request; can take a string, or any object which
// is treated as though it were JSON. TODO(nwittstock): support other things
func (r *Request) SetBody(b interface{}) error {
	switch b.(type) {
	case string:
		r.body = []byte(b.(string))
		return nil
	default:
		body, err := json.Marshal(b)
		if err != nil {
			return err
		}
		r.body = body
		return nil
	}

}

// Body retrieves the request body that's been set
func (r Request) Body() []byte {
	return r.body
}

// Do performs a Request and returns a Response
func (r Request) Do() (interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URI, nil)
	if err != nil {
		return nil, err
	}

	for h, v := range r.Headers {
		req.Header.Add(h, v)
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

	return body, nil
}
