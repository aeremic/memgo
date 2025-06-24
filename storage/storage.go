package storage

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

func (s *Storage) Get(k string) Value {
	if k == "" {
		return Value{}
	}

	result, ok := s.Map[Key(k)]
	if !ok {
		return Value{}
	}

	return result
}

func (s *Storage) Set(k string, d string) Result {
	s.Map[Key(k)] = Value{Data: d}

	return Result{true, 0}
}
