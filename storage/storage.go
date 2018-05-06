// Storage is package used for sharing various dynamic values across several modules
package storage

var store Storage

type Storage struct {
	vals     map[string]interface{}
	prefix   string
	onChange []func(Storage, string)
	subs     map[string]*Storage
}

func init() {
	store = New()
}

func Of(prefix string) *Storage {
	return store.Of(prefix)
}

func New() Storage {
	return Storage{map[string]interface{}{}, "", []func(Storage, string){}, make(map[string]*Storage)}
}

func (s *Storage) Of(prefix string) *Storage {
	if st, ok := s.subs[prefix]; ok {
		return st
	}
	onChange := make([]func(Storage, string), 0)
	st := Storage{s.vals, s.prefix + prefix + ".", onChange, make(map[string]*Storage)}
	s.subs[prefix] = &st
	return &st
}
func (s Storage) Set(key string, val interface{}) Storage {
	s.vals[s.prefix+key] = val
	s.changed(key)
	return s
}

func (s Storage) Bool(key string) bool {
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return false
	}
	return val.(bool)
}
func (s Storage) Int(key string) int {
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return 0
	}
	return val.(int)
}
func (s Storage) Float(key string) float64 {
	val, ok := s.vals[s.prefix+key]
	if !ok {
		return 0.0
	}
	return val.(float64)
}
func (s Storage) String(key string) string {
	v, ok := s.vals[s.prefix+key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (s *Storage) AddFloat(key string, diff float64) float64 {

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
	s.onChange = append(s.onChange, cb)
}

func (s Storage) changed(key string) {
	for _, f := range s.onChange {
		f(s, s.prefix+key)
	}
}
