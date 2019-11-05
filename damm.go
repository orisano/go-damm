package damm

import (
	"fmt"
)

type Alphabet interface {
	Chr(x uint8) byte
	Ord(c byte) (uint8, error)
	Size() uint8
}

var Digit digit
var UpperHex upperHex
var LowerHex lowerHex

func CheckSymbol(alphabet Alphabet, b []byte) (byte, error) {
	s, err := consume(alphabet, b)
	if err != nil {
		return 0, err
	}
	return alphabet.Chr(s), nil
}

func AppendCheckSymbol(alphabet Alphabet, b []byte) ([]byte, error) {
	sym, err := CheckSymbol(alphabet, b)
	if err != nil {
		return nil, err
	}
	return append(b, sym), nil
}

func IsValid(alphabet Alphabet, b []byte) (bool, error) {
	s, err := consume(alphabet, b)
	if err != nil {
		return false, err
	}
	return s == 0, nil
}

func consume(a Alphabet, b []byte) (uint8, error) {
	var s uint8
	sz := a.Size()
	mat, ok := matrices[sz]
	if !ok {
		return 0, fmt.Errorf("unsupported alphabet size: %d", sz)
	}
	for _, c := range b {
		o, err := a.Ord(c)
		if err != nil {
			return 0, err
		}
		s = mat[s][o]
	}
	return s, nil
}

type invalidCharacterError struct {
	regexp string
	c      byte
}

func (e *invalidCharacterError) Error() string {
	return fmt.Sprintf("got invalid character '0x%02x', must be %s", e.c, e.regexp)
}

type digit struct{}

func (digit) Chr(x uint8) byte {
	return '0' + x
}

func (digit) Ord(c byte) (uint8, error) {
	if c < '0' || '9' < c {
		return 0, &invalidCharacterError{"[0-9]", c}
	}
	return c - '0', nil
}

func (digit) Size() uint8 {
	return 10
}

type upperHex struct{}

func (upperHex) Chr(x uint8) byte {
	if x < 10 {
		return '0' + x
	} else {
		return 'A' + x - 10
	}
}

func (upperHex) Ord(c byte) (uint8, error) {
	if '0' <= c && c <= '9' {
		return c - '0', nil
	}
	if 'A' <= c && c <= 'F' {
		return c - 'A' + 10, nil
	}
	return 0, &invalidCharacterError{"[0-9A-F]", c}
}

func (upperHex) Size() uint8 {
	return 16
}

type lowerHex struct{}

func (lowerHex) Chr(x uint8) byte {
	if x < 10 {
		return '0' + x
	} else {
		return 'a' + x - 10
	}
}

func (lowerHex) Ord(c byte) (uint8, error) {
	if '0' <= c && c <= '9' {
		return c - '0', nil
	}
	if 'a' <= c && c <= 'f' {
		return c - 'a' + 10, nil
	}
	return 0, &invalidCharacterError{"[0-9a-f]", c}
}

func (lowerHex) Size() uint8 {
	return 16
}
