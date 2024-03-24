package chip8

import (
	"errors"
)

type Stack struct {
	data [32]int16
	sp int	// points to the next free element
}

func (s *Stack) push(pc int16) error {
	if s.sp > 31 {
		return errors.New("Stack overflow")
	}
	s.data[s.sp] = pc
	s.sp++
	return nil
}

func (s *Stack)pop() (int16, error) {
	var res int16
	if s.sp < 1 {
		return res, errors.New("Stack underflow")
	}
	res = s.data[s.sp]
	s.sp--
	return res, nil
}
