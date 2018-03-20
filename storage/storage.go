package storage

type Storage map[string]Node

func New () Storage {
	return Storage{}
}

func (s Storage) Get(key string) interface{} {
	val := s[key]
	switch v := val.(type) {
	case *Bool:
		return v.Get().(bool)
	case *Int:
		return v.Get().(int)
	case *String:
		return v.Get().(string)
	case *List:
		return v.Get().([]interface{})
	default:
		return v.Get()
	}
}
func (s Storage) Set(key string, val interface{}) {
	switch v := val.(type) {
	case int:
		intVal := Int(v)
		s[key] = &intVal
	case string:
		stringVal := String(v)
		s[key] = &stringVal
	case bool:
		boolVal := Bool(v)
		s[key] = &boolVal
	case []interface{}:
		listVal := List(v)
		s[key] = &listVal
	default: 
		panic("Unsupported type for storage val")
	}
}

type Node interface {
	Get() interface{}
	Set(interface{})
}

// types implementing Node interface

type Bool bool
func (b *Bool) Get() interface{} {
	return bool(*b)
}
func (b *Bool) Set(val interface{}) {
	*b = Bool(val.(bool))
}

type Int int
func (i *Int) Get() interface{} {
	return int(*i)
}
func (i *Int) Set(val interface{}) {
	*i = Int(val.(int))
}

type String string
func (s *String) Get() interface{} {
	return string(*s)
}
func (s *String) Set(val interface{}) {
	*s = String(val.(string))
}

type List []interface{}
func (l *List) Get() interface{} {
	return []interface{}(*l)
}
func (l *List) Set(val interface{}) {
	*l = []interface{}(val.([]interface{}))
}
