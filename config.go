package congruent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/imdario/mergo"
)

type Merger interface {
	Merge(interface{}) (interface{}, error)
}

// Headers represents HTTP header key/value pairs
type Headers map[string]string

// Authentication allows passing auth data which will be encoded later; this
// can be provided in Headers, but provided here we'll do the encoding for you
type Authentication struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a Authentication) String() string {
	return a.BasicAuth()
}

func (a Authentication) BasicAuth() string {
	s := fmt.Sprintf("%s:%s", a.Username, a.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))

	return fmt.Sprintf("Basic %s", encoded)
}

func (a Authentication) Merge(o *Authentication) (*Authentication, error) {
	if err := mergo.MergeWithOverwrite(a, o); err != nil {
		return nil, err
	}

	return &a, nil
}

// Server represents a server that will be requested against
type Server struct {
	Headers        Headers        `json:"headers"`
	Authentication Authentication `json:"authentication"`
	BaseURI        string         `json:"base_uri"`
}

func (s Server) Merge(o *Server) (*Server, error) {
	if err := mergo.MergeWithOverwrite(s, o); err != nil {
		return nil, err
	}

	return &s, nil
}

// Servers is an array of Server objects
type Servers []Server

// Request represents a request to be made
type Request struct {
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

func (r Request) Merge(o *Request) (*Request, error) {
	if err := mergo.MergeWithOverwrite(r, o); err != nil {
		return nil, err
	}

	return &r, nil
}

// Requests is an array of Request objects
type Requests []Request

// Config represents a whole configuration for a job
type Config struct {
	Global   Server   `json:"_global"`
	Servers  Servers  `json:"servers"`
	Requests Requests `json:"requests"`
}

// NewConfigFromJSON loads configuration data from a []byte of JSON data
func NewConfigFromJSON(data []byte) (*Config, error) {
	config := &Config{}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
