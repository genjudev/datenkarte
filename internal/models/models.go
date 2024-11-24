package models

import "datenkarte/internal/handlers"

// AuthHeader defines a single authentication header
type AuthHeader struct {
	Name  string `yaml:"name"`
	Match string `yaml:"match"`
}

// Auth defines the authentication headers
type Auth struct {
	Bearer string `yaml:"bearer"`
}

// HTTPHeader defines a single HTTP header
type HttpHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// AuthConfig defines the authentication configuration for HTTP requests
type AuthConfig struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

// Mapping defines the JSON mapping for a specific field
type Mapping struct {
	Name       string   `yaml:"name"`
	To         string   `yaml:"to"`
	Required   bool     `yaml:"required"`
	Nested     string   `yaml:"nested"`
	Fill       *Fill    `yaml:"fill"`
	InsertInto string   `yaml:"insert_into"`
	Handlers   []string `yaml:"handlers"`
}
type Fill struct {
	Type   string      `yaml:"type"`
	Value  interface{} `yaml:"value"`
	Prefix string      `yaml:"prefix,omitempty"`
}

type Validation struct {
	Field   string `yaml:"field"`
	Type    string `yaml:"type"`
	Pattern string `yaml:"pattern,omitempty"`
}

// EachLine defines the configuration for a single processing line
type EachLine struct {
	Map        []Mapping    `yaml:"map"`
	Validation []Validation `yaml:"validation"`
}

type HttpAuth struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type HttpType struct {
	Url        string       `yaml:"url"`
	Method     string       `yaml:"method"`
	Headers    []HttpHeader `yaml:"headers"`
	Auth       *HttpAuth    `yaml:"auth"`
	PayloadKey string       `yaml:"payload_key"`
}

// Rule defines the processing rules for an endpoint
type Rule struct {
	ID        string     `yaml:"id"`
	Delimiter string     `yaml:"delimiter"`
	Type      string     `yaml:"type"`
	Http      *HttpType  `yaml:"http"`
	EachLine  []EachLine `yaml:"each_line"`
}

// Config represents the entire YAML configuration
type Config struct {
	Auth     Auth               `yaml:"Auth"`
	Rules    []Rule             `yaml:"Rules"`
	Handlers []handlers.Handler `yaml:"Handlers"`
}
