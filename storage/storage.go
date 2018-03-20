package storage

type Storage  struct {
	vals map[string]interface{}
	prefix string
	onChange func(Storage, string)
}

func New () Storage {
	return Storage{map[string]interface{}{}, "", nil}
}

func (s Storage) Of(prefix string) Storage {
	return Storage{s.vals, s.prefix + prefix + ".", s.onChange}
}
func (s Storage) Set(key string, val interface{}) Storage {
	s.vals[s.prefix + key] = val
	s.changed(key)
	return s
}


func (s Storage) Bool(key string) bool {
	return s.vals[s.prefix+key].(bool)
}
func (s Storage) Int(key string) int {
	return s.vals[s.prefix+key].(int)
}
func (s Storage) Float(key string) float64 {
	return s.vals[s.prefix+key].(float64)
}
func (s Storage) String(key string) string {
	v, ok := s.vals[s.prefix+key]
	if !ok {
		return ""
	}
	return v.(string)
}

func (s *Storage) OnChange(cb func(Storage, string)) {
	(*s).onChange = cb
}

func (s Storage) changed(key string) {
	if s.onChange != nil {
		s.onChange(s, s.prefix+key)
	}
}