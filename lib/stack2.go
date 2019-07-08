package lib

import (
	"errors"
	"sync"
)

type stack struct {
	lock sync.Mutex // you don't have to do this if you don't want thread safety
	s    []int
}

func New() *stack {
	return &stack{sync.Mutex{}, make([]int, 0)}
}

func (s *stack) Push(v int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.s = append(s.s, v)
}

func (s *stack) Pop() (int, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	l := len(s.s)
	if l == 0 {
		return 0, errors.New("empty Stack")
	}
	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}
