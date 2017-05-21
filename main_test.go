package main

import "testing"

func TestCalculateWidths(t *testing.T) {
	size := 140
	boldWidth, width, err := CalculateWidths(&size)

	if boldWidth != 18. {
		t.Error("Invalid calculation of boldWidth")
	}

	if width != 1.8 {
		t.Error("Invalid calculation of width")
	}

	if err != nil {
		t.Error("Error happened")
	}
}

func TestCalculateWidthsSizeOverConversionFactor(t *testing.T) {
	size := int(ConversionFactor) * 10
	_, _, err := CalculateWidths(&size)

	if err == nil {
		t.Error("boldWidth is equal 0 or negative")
	}
}

func TestCalculateLineData(t *testing.T) {
	boldWidth := 20.
	width := 2.
	data, err := CalculateLineData(boldWidth, width)

	if err != nil {
		t.Error("no error should be happen")
	}

	if data.BoldWidth != boldWidth {
		t.Error("boldWidth in struct is missmatch original value")
	}

	if data.Width != width {
		t.Error("width in struct is missmatch original value")
	}

	if data.HorizontalMargin != 8.5 {
		t.Error("width in struct is missmatch original value")
	}

	if data.VerticalMargin != 5. {
		t.Error("width in struct is missmatch original value")
	}

	if data.MaxHeight != 288.5 {
		t.Error("width in struct is missmatch original value")
	}

	if data.MaxWidth != 205. {
		t.Error("width in struct is missmatch original value")
	}
}

