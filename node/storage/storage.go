package storage

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"sync"

	common "memgo_common"

	"github.com/xyproto/jpath"
)

type Key string
type Value struct {
	Data any
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

func (s *Storage) Dispose() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s = nil

	return common.SUCCESS
}

func (s *Storage) Get(k string) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	if k == "" {
		return common.EMPTY_KEY
	}

	result, ok := s.Map[Key(k)]
	if !ok {
		return common.NOT_FOUND
	}

	return result.Data
}

func (s *Storage) GetAll() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out bytes.Buffer

	elements := []string{}
	for key, value := range s.Map {
		switch v := value.Data.(type) {
		case int:
			elements = append(elements, fmt.Sprintf(`"%s":%d`, key, v))
		case float64:
			elements = append(elements, fmt.Sprintf(`"%s":%f`, key, v))
		case bool:
			elements = append(elements, fmt.Sprintf(`"%s":%t`, key, v))
		case string:
			elements = append(elements, fmt.Sprintf(`"%s":"%s"`, key, v))
		default:
			elements = append(elements, fmt.Sprintf(`"%s":%v`, key, v))
		}
	}

	if len(elements) > 0 {
		out.WriteString("{")
		out.WriteString(strings.Join(elements, ", "))
		out.WriteString("}")

		return out.String()
	}

	return common.NOT_FOUND
}

func (s *Storage) Set(k string, d string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var parsed any

	if integer, err := strconv.Atoi(d); err == nil {
		parsed = integer
	} else if float, err := strconv.ParseFloat(d, 64); err == nil {
		parsed = float
	} else if boolean, err := strconv.ParseBool(d); err == nil {
		parsed = boolean
	} else {
		parsed = d
	}

	s.Map[Key(k)] = Value{Data: parsed}

	return common.SUCCESS
}

func (s *Storage) Delete(k string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.Map, Key(k))

	return common.SUCCESS
}

func (s *Storage) DeleteAll() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Map = make(map[Key]Value)

	return common.SUCCESS
}

func (s *Storage) GetByKeyAndPath(k string, p string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := s.Get(k)

	strRes, ok := res.(string)
	if !ok {
		return common.ERROR
	}

	document, err := jpath.New([]byte(strRes))
	if err != nil {
		return err.Error()
	}

	node := document.GetNode(p)
	if node != nil && node.String() != "" {
		return node.String()
	}

	return common.NOT_FOUND
}

func (s *Storage) SelectByPath(p string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var out bytes.Buffer

	elements := []string{}
	for key, value := range s.Map {
		strRes, ok := value.Data.(string)
		if !ok {
			continue
		}

		document, err := jpath.New([]byte(strRes))
		if err != nil {
			continue
		}

		node := document.GetNode(p)
		if node != nil && node.String() != "" {
			elements = append(elements, fmt.Sprintf(`"%s":%s`, key, strRes))
		}
	}

	if len(elements) > 0 {
		out.WriteString("{")
		out.WriteString(strings.Join(elements, ", "))
		out.WriteString("}")

		return out.String()
	}

	return common.NOT_FOUND
}
