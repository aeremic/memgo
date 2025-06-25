package storage

import (
	"bytes"
	"fmt"
	"strings"
)

type Result struct {
	Success bool
	Error   int
}

const (
	_ int = iota
	EMPTY_KEY
)

type Key string
type Value struct {
	Data string
}

type Storage struct {
	Map map[Key]Value
}

func New() *Storage {
	s := &Storage{}
	s.Map = make(map[Key]Value)

	return s
}

func (s *Storage) Get(k string) string {
	if k == "" {
		return ""
	}

	result, ok := s.Map[Key(k)]
	if !ok {
		return ""
	}

	return result.Data
}

func (s *Storage) GetAll() string {
	var out bytes.Buffer

	elements := []string{}
	for key, value := range s.Map {
		elements = append(elements, fmt.Sprintf(`"%s":"%s"`, key, value.Data))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}

func (s *Storage) Set(k string, d string) Result {
	s.Map[Key(k)] = Value{Data: d}

	return Result{true, 0}
}
