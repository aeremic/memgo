package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	Url   string       `json:"url"`
	Port  string       `json:"port"`
	Nodes []ConfigNode `json:"nodes"`
}

type ConfigNode struct {
	Url  string `json:"url"`
	Port string `json:"port"`
}

func Get(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := file
	decoder := json.NewDecoder(reader)

	var c Config
	decoderErr := decoder.Decode(&c)
	if decoderErr != nil {
		return nil, decoderErr
	}

	return &c, nil
}
