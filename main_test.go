package main

import "testing"

func TestCalculateVerticalMargins(t *testing.T) {
	if margin := CalculateVerticalMargins(20.); margin != MARGIN_VERTICAL {
		t.Error("Invalid margin")
	}
}
