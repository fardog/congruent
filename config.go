package congruent

import (
	"encoding/json"
)

// Headers represents HTTP header key/value pairs
type Headers map[string]string

// Authentication allows passing auth data which will be encoded later; this
// can be provided in Headers, but provided here we'll do the encoding for you
type Authentication struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Server represents a server that will be requested against
type Server struct {
	Headers        Headers        `json:"headers"`
	Authentication Authentication `json:"authentication"`
	BaseURI        string         `json:"base_uri"`
}

// Servers is an array of Server objects
type Servers []Server

// Expectation represents an expected outcome from a request
type Expectation struct {
	Headers Headers     `json:"headers"`
	Status  int32       `json:"status"`
	Body    interface{} `json:"body"`
}

// Request represents a request to be made
type Request struct {
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
	Expect Expectation `json:"expect"`
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
