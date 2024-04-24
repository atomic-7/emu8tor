package chip8

import (
	"errors"
)

type Stack struct {
	data [32]uint16
	sp int	// points to the next free element
}

func (s *Stack) push(pc uint16) error {
	if s.sp > 31 {
		return errors.New("C8: Stack overflow")
	}
	s.data[s.sp] = pc
	s.sp++
	return nil
}

func (s *Stack)pop() (uint16, error) {
	var res uint16
	if s.sp < 1 {
		return res, errors.New("C8: Stack underflow")
	}
	res = s.data[s.sp - 1]
	s.sp--
	return res, nil
}
