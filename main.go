package main

import (
	"errors"
	"flag"
	"github.com/jung-kurt/gofpdf"
	"math"
	"strconv"
)

// 120 20mm, 140 18mm, 160 16mm, 180 14mm, 200 12mm
const (
	Width            = 210
	Height           = 297
	MarginVertical   = 5.
	MarginHorizontal = 8.
	LineWidthBold    = 0.5
	LineWidth        = 0.2
	ConversionFactor = 32.
	Vertical         = iota
	Horizontal
)

type LineCalc struct {
	BoldWidth        float64
	Width            float64
	VerticalMargin   float64
	HorizontalMargin float64
	MaxWidth         float64
	MaxHeight        float64
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
	filename, size := getFlags()
	boldWidth, width, err := CalculateWidths(size)

	if err != nil {
		panic(err.Error())
	}

	pdf := initPdf()
	data, err := CalculateLineData(boldWidth, width)
	pdf.Text(data.VerticalMargin, data.HorizontalMargin-2., strconv.Itoa(*size))

	if err != nil {
		panic(err.Error())
	}

	DrawVerticalLines(pdf, data)
	DrawHorizontalLines(pdf, data)

	err = pdf.OutputFileAndClose(*filename)

	if err != nil {
		panic(err.Error())
	}
}

// Calculates box widths.
// Each bold quad has 10 smaller quads.
// The size of bold quad is set by the type and it converts from the
// carpet size to millimeter
func CalculateWidths(size *int) (float64, float64, error) {
	boldWidth := truncate(ConversionFactor - float64(*size)/10.)

	if boldWidth <= 0 {
		return 0., 0., errors.New("invalid type")
	}

	width := boldWidth / 10.

	return boldWidth, width, nil
}

// CalculateLineData calculate the vertical and horizontal margins.
// It also calculates the maximum with and height from the resulting margins
func CalculateLineData(boldWidth, width float64) (*LineCalc, error) {
	verticalMargin, err := calculateMargin(boldWidth, Width, MarginVertical)

	if err != nil {
		return nil, err
	}

	horizontalMargin, err := calculateMargin(boldWidth, Height, MarginHorizontal)

	if err != nil {
		return nil, err
	}

	return &LineCalc{
		BoldWidth:        boldWidth,
		Width:            width,
		VerticalMargin:   verticalMargin,
		HorizontalMargin: horizontalMargin,
		MaxWidth:         Width - verticalMargin,
		MaxHeight:        Height - horizontalMargin,
	}, nil
}

// DrawVerticalLines get the calculated stuff and draw all the vertical lines
func DrawVerticalLines(pdf *gofpdf.Fpdf, data *LineCalc) {
	drawLines(pdf, Vertical, &Line{
		BoldWidth:  data.BoldWidth,
		Width:      data.Width,
		InitMargin: data.VerticalMargin,
		LineMargin: data.HorizontalMargin,
		Max:        data.MaxWidth,
		LineMax:    data.MaxHeight,
	})
}

// DrawHorizontalLines get the calculated stuff and draw all the horizontal lines
func DrawHorizontalLines(pdf *gofpdf.Fpdf, data *LineCalc) {
	drawLines(pdf, Horizontal, &Line{
		BoldWidth:  data.BoldWidth,
		Width:      data.Width,
		InitMargin: data.HorizontalMargin,
		LineMargin: data.VerticalMargin,
		Max:        data.MaxHeight,
		LineMax:    data.MaxWidth,
	})
}

func drawLines(pdf *gofpdf.Fpdf, lineType int, data *Line) error {
	j := 0.

	for i := data.InitMargin; i <= data.Max; i = truncate(i + data.Width) {
		pdf.SetLineWidth(LineWidth)

		if math.Mod(j, 10) == 0 {
			pdf.SetLineWidth(LineWidthBold)
		}

		switch lineType {
		case Vertical:
			pdf.Line(i, data.LineMargin, i, data.LineMax)
		case Horizontal:
			pdf.Line(data.LineMargin, i, data.LineMax, i)
		default:
			return errors.New("unknown type")
		}
		j++
	}

	return nil
}

func calculateMargin(width, total, margin float64) (float64, error) {
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

func truncate(some float64) float64 {
	return float64(int(some*100)) / 100
}

func getFlags() (*string, *int) {
	filename := flag.String("filename", "carpet.pdf", "output filename, default: carpet.pdf")
	size := flag.Int("type", 120, "carpet type")
	flag.Parse()

	return filename, size
}

func initPdf() *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetDrawColor(0, 0, 0)
	pdf.SetFont("Helvetica", "", 8)

	return pdf
}
