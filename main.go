package main

import (
	"errors"
	"github.com/jung-kurt/gofpdf"
	"math"
)

const (
	WIDTH             = 210
	HEIGHT            = 297
	MARGIN_VERTICAL   = 5.
	MARGIN_HORIZONTAL = 8.
	LINE_BOLD         = 0.5
	LINE              = 0.2
	VERTICAL          = iota
	HORIZONTAL
)

type LineCalc struct {
	BoldWidth        float64
	Width            float64
	VerticalMargin   float64
	HorizontalMargin float64
	MaxWidth         float64
	maxHeight        float64
}

type Line struct {
	BoldWidth  float64
	Width      float64
	InitMargin float64
	LineMargin float64
	Max        float64
	LineMax    float64
}

func main() {
	boldWidth := 20.
	width := 2.

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetDrawColor(0, 0, 0)

	data, err := CalculateLineData(boldWidth, width)

	if err != nil {
		panic(err.Error())
	}

	DrawVerticalLines(pdf, data)
	DrawHorizontalLines(pdf, data)

	err = pdf.OutputFileAndClose("test1.pdf")

	if err != nil {
		panic(err.Error())
	}
}

func CalculateLineData(boldWidth, width float64) (*LineCalc, error) {
	verticalMargin, err := CalculateMargin(boldWidth, WIDTH, MARGIN_VERTICAL)

	if err != nil {
		return nil, err
	}

	horizontalMargin, err := CalculateMargin(boldWidth, HEIGHT, MARGIN_HORIZONTAL)

	if err != nil {
		return nil, err
	}

	return &LineCalc{
		BoldWidth:        boldWidth,
		Width:            width,
		VerticalMargin:   verticalMargin,
		HorizontalMargin: horizontalMargin,
		MaxWidth:         WIDTH - verticalMargin,
		maxHeight:        HEIGHT - horizontalMargin,
	}, nil
}

func DrawVerticalLines(pdf *gofpdf.Fpdf, data *LineCalc) {
	DrawLines(pdf, VERTICAL, &Line{
		BoldWidth:  data.BoldWidth,
		Width:      data.Width,
		InitMargin: data.VerticalMargin,
		LineMargin: data.HorizontalMargin,
		Max:        data.MaxWidth,
		LineMax:    data.maxHeight,
	})
}

func DrawHorizontalLines(pdf *gofpdf.Fpdf, data *LineCalc) {
	DrawLines(pdf, HORIZONTAL, &Line{
		BoldWidth:  data.BoldWidth,
		Width:      data.Width,
		InitMargin: data.HorizontalMargin,
		LineMargin: data.VerticalMargin,
		Max:        data.maxHeight,
		LineMax:    data.MaxWidth,
	})
}

func DrawLines(pdf *gofpdf.Fpdf, lineType int, data *Line) error {
	j := 0.

	for i := data.InitMargin; i <= data.Max; i += data.Width {
		pdf.SetLineWidth(LINE)

		if math.Mod(j, data.BoldWidth) == 0 {
			pdf.SetLineWidth(LINE_BOLD)
		}

		switch lineType {
		case VERTICAL:
			pdf.Line(i, data.LineMargin, i, data.LineMax)
		case HORIZONTAL:
			pdf.Line(data.LineMargin, i, data.LineMax, i)
		default:
			return errors.New("unknown type")
		}
		j += data.Width
	}

	return nil
}

func CalculateMargin(width, total, margin float64) (float64, error) {
	rest := math.Mod(total-(2*margin), width)

	if rest == 0 {
		return margin, nil
	}

	newMargin := margin + (rest / 2)

	if rest = math.Mod(total-(2*newMargin), width); rest == 0 {
		return newMargin, nil
	}

	return 0., errors.New("Invalid Margin")
}
