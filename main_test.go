package main

import "testing"

func TestCalculateVerticalMargins(t *testing.T) {
	if margin := CalculateVerticalMargins(20.); margin != MarginVertical {
		t.Error("Invalid margin")
	}
}
