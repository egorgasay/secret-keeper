package config

import (
	"flag"
	"secret-keeper/internal/server/storage"
)

type Config struct {
	Host     string
	DBConfig storage.Config
}

// Flag struct for parsing from env and cmd args.
type Flag struct {
	Host *string `json:"server_address,omitempty"`
	URI  *string `json:"uri,omitempty"`
}

var f Flag

func init() {
	f.Host = flag.String("a", defaults["Host"], "-a=host")
	f.URI = flag.String("u", defaults["URI"], "-u=uri")
}

const (
	defaultHost = "127.0.0.1:8080"
	defaultURI  = "127.0.0.1:800"
)

var defaults = map[string]string{
	"Host": defaultHost,
	"URI":  defaultURI,
}

func New() (*Config, error) {
	flag.Parse()

	return &Config{
		Host: *f.Host,
		DBConfig: storage.Config{
			URI: *f.URI,
		},
	}, nil
}
