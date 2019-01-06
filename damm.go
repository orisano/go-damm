package damm

import (
	"fmt"
	"io"
)

type Status interface {
	io.Writer
	IsValid() bool
	CheckSymbol() byte
	Reset()
}

type Alphabet interface {
	Chr(x uint8) byte
	Ord(c byte) (uint8, error)
	Size() uint8
}

func New(alphabet Alphabet) (Status, error) {
	sz := alphabet.Size()
	mat, ok := matrices[sz]
	if !ok {
		return nil, fmt.Errorf("unsupported alphabet size: %d", sz)
	}
	return &status{0, alphabet, mat}, nil
}

func Digit() Alphabet {
	return &digit{}
}

func UpperHex() Alphabet {
	return &upperHex{}
}

func LowerHex() Alphabet {
	return &lowerHex{}
}

func CheckSymbol(b []byte, alphabet Alphabet) (byte, error) {
	st, err := New(alphabet)
	if err != nil {
		return 0, err
	}
	if _, err := st.Write(b); err != nil {
		return 0, err
	}
	return st.CheckSymbol(), nil
}

func CheckDigit(b []byte) (byte, error) {
	return CheckSymbol(b, &digit{})
}

func CheckUpperHex(b []byte) (byte, error) {
	return CheckSymbol(b, &upperHex{})
}

func CheckLowerHex(b []byte) (byte, error) {
	return CheckSymbol(b, &lowerHex{})
}

func IsValid(b []byte, alphabet Alphabet) (bool, error) {
	st, err := New(alphabet)
	if err != nil {
		return false, err
	}
	if _, err := st.Write(b); err != nil {
		return false, err
	}
	return st.IsValid(), nil
}

func IsValidDigit(b []byte) (bool, error) {
	return IsValid(b, &digit{})
}

func IsValidUpperHex(b []byte) (bool, error) {
	return IsValid(b, &upperHex{})
}

func IsValidLowerHex(b []byte) (bool, error) {
	return IsValid(b, &lowerHex{})
}

type status struct {
	s        uint8
	alphabet Alphabet
	mat      [][]uint8
}

func (s *status) Write(p []byte) (n int, err error) {
	for _, c := range p {
		o, err := s.alphabet.Ord(c)
		if err != nil {
			return 0, err
		}
		s.s = s.mat[s.s][o]
	}
	return len(p), nil
}

func (s *status) IsValid() bool {
	return s.s == 0
}

func (s *status) CheckSymbol() byte {
	return s.alphabet.Chr(s.s)
}

func (s *status) Reset() {
	s.s = 0
}

type invalidCharacterError struct {
	regexp string
	c      byte
}

func (e *invalidCharacterError) Error() string {
	return fmt.Sprintf("got invalid character '0x%02x', must be %s", e.c, e.regexp)
}

type digit struct{}

func (*digit) Chr(x uint8) byte {
	return '0' + x
}

func (*digit) Ord(c byte) (uint8, error) {
	if c < '0' || '9' < c {
		return 0, &invalidCharacterError{"[0-9]", c}
	}
	return c - '0', nil
}

func (*digit) Size() uint8 {
	return 10
}

type upperHex struct{}

func (*upperHex) Chr(x uint8) byte {
	if x < 10 {
		return '0' + x
	} else {
		return 'A' + x - 10
	}
}

func (*upperHex) Ord(c byte) (uint8, error) {
	if '0' <= c && c <= '9' {
		return c - '0', nil
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A' + 10, nil
	}
	return 0, &invalidCharacterError{"[0-9A-F]", c}
}

func (*upperHex) Size() uint8 {
	return 16
}

type lowerHex struct{}

func (*lowerHex) Chr(x uint8) byte {
	if x < 10 {
		return '0' + x
	} else {
		return 'a' + x - 10
	}
}

func (*lowerHex) Ord(c byte) (uint8, error) {
	if '0' <= c && c <= '9' {
		return c - '0', nil
	}
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10, nil
	}
	return 0, &invalidCharacterError{"[0-9a-f]", c}
}

func (*lowerHex) Size() uint8 {
	return 16
}
