package ssnrgo

import (
	"errors"
)

type Config struct {
	Port string
	Host string
	Name string
}

func NewConfig(host, port, name string) (*Config, error) {
	r := new(Config)
	if host == "" {
		return nil, errors.New("Missing host address")
	}
	if port == "" {
		return nil, errors.New("Missing application's port")
	}
	r.Port = port
	r.Host = host
	r.Name = name
	return r, nil
}
