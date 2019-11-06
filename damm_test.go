package damm_test

import (
	"testing"

	"github.com/orisano/go-damm"
)

func TestIsValid(t *testing.T) {
	ts := []string{
		"0123456789",
		"11223334442222",
		"9876",
		"100101010102030",
	}
	as := []damm.Alphabet{
		damm.Digit,
		damm.LowerHex,
		damm.UpperHex,
	}
	for _, a := range as {
		for _, tc := range ts {
			testIsValid(t, a, []byte(tc))
		}
	}
}

func testIsValid(t *testing.T, a damm.Alphabet, b []byte) {
	t.Helper()
	d, err := damm.CheckSymbol(a, b)
	if err != nil {
		t.Errorf("failed to calculate check digit: %v", err)
		return
	}
	b = append(b, d)
	for i := 0; i < len(b)-1; i++ {
		if b[i] == b[i+1] {
			continue
		}
		var bb []byte
		bb = append(bb, b...)
		bb[i], bb[i+1] = bb[i+1], bb[i]
		ok, err := damm.IsValid(a, bb)
		if err != nil {
			t.Errorf("failed to check check digit: %v", err)
			continue
		}
		if ok {
			t.Errorf("unexpected result of validation, expected to invalid: %q", string(bb))
		}
	}
	for i := 0; i < len(b); i++ {
		for ord := uint8(0); ord < a.Size(); ord++ {
			c := a.Chr(ord)
			var bb []byte
			bb = append(bb, b...)
			bb[i] = c
			ok, err := damm.IsValid(a, bb)
			if err != nil {
				t.Errorf("failed to check check digit: %v", err)
				continue
			}
			if (c == b[i]) != ok {
				t.Errorf("unexpected result of validation, expected to invalid: %q", string(bb))
			}
		}
	}
}
