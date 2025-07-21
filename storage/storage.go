package storage

import (
	"bytes"
	"fmt"
	"strings"
	"sync"

	"github.com/xyproto/jpath"
)

const (
	EMPTY_KEY = "EMPTY_KEY"
	NOT_FOUND = "NOT_FOUND"
	EMPTY     = "EMPTY"
	SUCCESS   = "SUCCESS"
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
		return EMPTY_KEY
	}

	result, ok := s.Map[Key(k)]
	if !ok {
		return NOT_FOUND
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

func (s *Storage) Set(k string, d string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Map[Key(k)] = Value{Data: d}

	return SUCCESS
}

func (s *Storage) Delete(k string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Map, Key(k))

	return SUCCESS
}

func (s *Storage) DeleteAll() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Map = make(map[Key]Value)

	return SUCCESS
}

func (s *Storage) GetByKeyAndPath(k string, p string) string {
	res := s.Get(k)

	document, err := jpath.New([]byte(res))
	if err != nil {
		return err.Error()
	}

	node := document.GetNode(p)
	if node != nil {
		node.String()
	}

	return NOT_FOUND
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

	if len(elements) > 0 {
		out.WriteString("{")
		out.WriteString(strings.Join(elements, ", "))
		out.WriteString("}")

		return out.String()
	}

	return NOT_FOUND
}
