package tests

import (
	"mamela/audio"
	"testing"
)

func PadTestSingle(t *testing.T) {
	padded := audio.Pad(0)
	if padded != "00" {
		t.Errorf("pad(0) = %s; want 00", padded)
	}
}
func PadTestDouble(t *testing.T) {
	padded := audio.Pad(11)
	if padded != "11" {
		t.Errorf("pad(11) = %s; want 11", padded)
	}
}
