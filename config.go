package ssnrgo

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port string `json:"port"`
	Host string `json:"host"`
	Name string `json:"name"`
}

func NewConfig(filePath string) *Config {
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		panic("Failed to open config file:" + filePath)
	}
	decoder := json.NewDecoder(f)
	r := new(Config)
	err = decoder.Decode(r)
	if err != nil {
		panic("Failed to decode file")
	}
	return r
}
