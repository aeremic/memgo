package storage

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/xyproto/jpath"
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
	mu  sync.Mutex
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
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Map[Key(k)] = Value{Data: d}

	return Result{true, 0}
}

func (s *Storage) Delete(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Map, Key(k))
}

func (s *Storage) DeleteAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Map = make(map[Key]Value)
}

func (s *Storage) GetByKeyAndPath(k string, p string) string {
	res := s.Get(k)

	document, err := jpath.New([]byte(res))
	if err != nil {
		return err.Error()
	}

	return document.GetNode(p).String()
}

func (s *Storage) Select(p string) string {
	var out bytes.Buffer

	elements := []string{}
	for key, value := range s.Map {
		document, err := jpath.New([]byte(value.Data))
		if err != nil {
			continue
		}

		_, exist := document.CheckGet(p)
		if exist {
			elements = append(elements, fmt.Sprintf(`"%s":"%s"`, key, value.Data))
		}
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")

	return out.String()
}
