package storage

import (
	"sync"
)

var store Storage

type Storage struct {
	vals     map[string]interface{}
	prefix   string
	onChange func(Storage, string)
	subs     map[string]*Storage
	mutex    *sync.Mutex
}

func init() {
	store = New()
}

func Of(prefix string) *Storage {
	return store.Of(prefix)
}

func New() Storage {
	return Storage{map[string]interface{}{}, "", nil, make(map[string]*Storage), &sync.Mutex{}}
}

func (s *Storage) Of(prefix string) *Storage {
	if st, ok := s.subs[prefix]; ok {
		return st
	}
	st := Storage{s.vals, s.prefix + prefix + ".", s.onChange, make(map[string]*Storage), s.mutex}
	s.subs[prefix] = &st
	return &st
}
func (s Storage) Set(key string, val interface{}) Storage {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.vals[s.prefix+key] = val
	s.changed(key)
	return s
}

func (s Storage) Bool(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return false
	}
	return val.(bool)
}
func (s Storage) Int(key string) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return 0
	}
	return val.(int)
}
func (s Storage) Float(key string) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return 0.0
	}
	return val.(float64)
}
func (s Storage) String(key string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	v, ok := s.vals[s.prefix+key]
	if !ok {
		return ""
	}
	return v.(string)
}


func (s *Storage) AddFloat(key string, diff float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	val, ok := s.vals[s.prefix+key].(float64)
	if !ok {
		val = 0.0
	}
	val += diff
	s.vals[s.prefix+key] = val
	s.changed(key)
	return val
}



func (s *Storage) OnChange(cb func(Storage, string)) {
	(*s).onChange = cb
}

func (s Storage) changed(key string) {
	if s.onChange != nil {
		s.onChange(s, s.prefix+key)
	}
}
