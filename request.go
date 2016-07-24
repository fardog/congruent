package congruent

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func NewRequest(m, uri string, h HeadersDef) *Request {
	return &Request{m, uri, h, nil}
}

type Request struct {
	Method  string
	URI     string
	Headers HeadersDef
	body    []byte
}

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

func (r Request) Body() []byte {
	return r.body
}

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
