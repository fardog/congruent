package congruent

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// HeadersDef represents HTTP header key/value pairs
type HeadersDef map[string]string

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

// ServerDef represents a server that will be requested against
type ServerDef struct {
	Headers        HeadersDef `json:"headers"`
	Authentication AuthDef    `json:"authentication"`
	BaseURI        string     `json:"base_uri"`
}

// ServerDefs is an array of Server objects
type ServerDefs []ServerDef

// RequestDef represents a request to be made
type RequestDef struct {
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Body   interface{} `json:"body"`
}

// RequestDefs is an array of Request objects
type RequestDefs []RequestDef

// Config represents a whole configuration for a job
type Config struct {
	Global   ServerDef   `json:"_global"`
	Servers  ServerDefs  `json:"servers"`
	Requests RequestDefs `json:"requests"`
}

// NewConfigFromJSON loads configuration data from a []byte of JSON data
func NewConfigFromJSON(data []byte) (*Config, error) {
	config := &Config{}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
