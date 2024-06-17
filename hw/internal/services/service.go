package services

import (
	my_errors "hw/internal/errors"
	"strconv"
	"sync"
)

var Mutex sync.Mutex

type Set struct {
	array []int
}

func NewSet() Set {
	return Set{
		array: make([]int, 0),
	}
}

func (s *Set) Add(url_elem string) error {
	element, err := strconv.Atoi(url_elem)
	if err != nil {
		return err
	}
	Mutex.Lock()
	for i := 0; i < len(s.array); i++ {
		if s.array[i] == element {
			return my_errors.ErrRepeat
		}
	}
	s.array = append(s.array, element)
	Mutex.Unlock()
	return nil
}

func (s *Set) DeleteElem(url_elem string) error {
	element, err := strconv.Atoi(url_elem)
	if err != nil {
		return err
	}
	Mutex.Lock()
	for idx := 0; idx < len(s.array); idx++ {
		if s.array[idx] == element {
			s.array = append(s.array[0:idx], s.array[idx+1:]...)
			return nil
		}
	}
	Mutex.Unlock()
	return my_errors.ErrNoElem
}

func (s *Set) DeleteAll() {
	Mutex.Lock()
	s.array = make([]int, 0)
	Mutex.Unlock()
}

func (s *Set) Check(url_element string) error {
	element, err := strconv.Atoi(url_element)
	if err != nil {
		return err
	}
	check := false
	Mutex.Lock()
	for idx := 0; idx < len(s.array); idx++ {
		if (s.array)[idx] == element {
			check = true
			return nil
		}
	}
	Mutex.Unlock()
	if !check {
		return my_errors.ErrNoElem
	}
	return nil
}
