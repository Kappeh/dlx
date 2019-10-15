package dlx

import "errors"

type Matrix struct {
	root             element
	headers, details []element
	rows             []row
	rowCount         int
	solution         solution
}

type element struct {
	size                          int
	row                           *row
	header, up, down, left, right *element
}

type row struct {
	index        int
	firstElement *element
	covered      bool
}

type solution struct {
	values   []int
	stackptr int
}

func (s *solution) push(i int) {
	if s.stackptr == len(s.values) {
		s.values = append(s.values, i)
	} else {
		s.values[s.stackptr] = i
	}
	s.stackptr++
}

func (s *solution) pop() (int, error) {
	if s.stackptr == 0 {
		return 0, errors.New("cannot perform pop: stack is empty")
	}
	s.stackptr--
	return s.values[s.stackptr], nil
}

func (s *solution) size() int {
	return s.stackptr
}
