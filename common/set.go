package common

import "fmt"

type Set struct {
	container map[interface{}]bool
}

func NewSet() *Set {
	return &Set{
		container: make(map[interface{}]bool),
	}
}

func (s *Set) Insert(v interface{}) bool {
	if _, ok := s.container[v]; ok {
		return false
	}
	s.container[v] = true
	return true
}

func (s *Set) Erase(v interface{}) {
	delete(s.container, v)
}

func (s *Set) AsStringArray() []string {
	rst := make([]string, 0)
	for k, _ := range s.container {
		v := fmt.Sprintf("%v", k)
		rst = append(rst, v)
	}

	return rst
}

func (s *Set) Exist(v interface{}) bool {
	if _, ok := s.container[v]; ok {
		return true
	}

	return false
}
