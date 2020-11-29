package chart

import (
	"fmt"
	"io"
	"math"
	"strings"

	svg "github.com/ajstarks/svgo"
)

const width = 800
const labelAreaWidth = 120
const titleHeight = 50
const barHeight = 20
const paddingBetweenBars = 2
const paddingBottom = 20
const paddingRight = 40
const textColor = "#1b1e23"

// BarChart implements the Chart interface.
type BarChart struct {
	Title     string
	Bars      []Bar
	Precision int
	canvas    *svg.SVG
}

// Bar represents a single bar in a bar chart.
type Bar struct {
	Value float64
	Label string
	Color string
}

// Render will write an svg into the given target.
func (chart *BarChart) Render(target io.Writer) {
	chart.canvas = svg.New(target)
	chart.renderBackground()
	chart.renderTitle()
	chart.renderBars()
	chart.renderYAxis()
	chart.canvas.End()
}

func (chart *BarChart) renderBackground() {
	height := titleHeight + len(chart.Bars)*(barHeight+paddingBetweenBars) + paddingBottom
	chart.canvas.Start(
		width,
		height,
		"font-family=\"-apple-system,BlinkMacSystemFont,Segoe UI,Helvetica,Arial,sans-serif,Apple Color Emoji,Segoe UI Emoji\"",
	)
	chart.canvas.Rect(0, 0, width, height, "fill=\"white\"")
}

func (chart *BarChart) renderTitle() {
	chart.canvas.Text(labelAreaWidth,
		titleHeight/2,
		chart.Title,
		"dominant-baseline=\"central\"",
		"font-size=\"1.2em\"",
		fmt.Sprintf("fill=\"%s\"", textColor),
	)
}

func (chart *BarChart) renderYAxis() {
	chart.canvas.Line(
		labelAreaWidth,
		titleHeight-paddingBetweenBars,
		labelAreaWidth,
		titleHeight+len(chart.Bars)*(barHeight+paddingBetweenBars),
		fmt.Sprintf("stroke=\"%s\"", textColor),
	)
}

func (chart *BarChart) renderBars() {
	for index, bar := range chart.Bars {
		chart.renderBar(index, &bar)
	}
}

func (chart *BarChart) renderBar(index int, bar *Bar) {
	barWidth := chart.scaleBar(bar.Value)
	offsetY := titleHeight + index*(barHeight+paddingBetweenBars)
	chart.canvas.Rect(
		labelAreaWidth,
		offsetY,
		barWidth,
		barHeight,
		"fill=\"steelblue\"",
	)
	chart.renderBarValue(bar.Value, barWidth, offsetY)
	chart.renderBarLabel(bar.Label, barWidth, offsetY)
}

func (chart *BarChart) renderBarValue(value float64, barWidth, offsetY int) {
	chart.canvas.Text(
		barWidth+labelAreaWidth-2*paddingBetweenBars,
		offsetY+barHeight/2,
		fmt.Sprintf("%s", chart.formatToPrecision(value)),
		"text-anchor=\"end\"",
		"dominant-baseline=\"central\"",
		"fill=\"white\"",
	)
}

func (chart *BarChart) renderBarLabel(labelText string, barWidth, offsetY int) {
	chart.canvas.Text(
		labelAreaWidth-10,
		offsetY+barHeight/2,
		labelText,
		"text-anchor=\"end\"",
		"dominant-baseline=\"central\"",
		fmt.Sprintf("fill=\"%s\"", textColor),
	)

	chart.canvas.Line(
		labelAreaWidth-7,
		offsetY+barHeight/2,
		labelAreaWidth,
		offsetY+barHeight/2,
		fmt.Sprintf("stroke=\"%s\"", textColor),
	)
}

func (chart *BarChart) scaleBar(barValue float64) int {
	ratio := barValue / chart.maxValue()
	maxWidth := float64(width - labelAreaWidth - paddingRight)
	return int(math.RoundToEven(ratio * maxWidth))
}

func (chart *BarChart) maxValue() float64 {
	maxValue := 0.0
	for _, bar := range chart.Bars {
		if bar.Value > maxValue {
			maxValue = bar.Value
		}
	}
	return maxValue
}

func (chart *BarChart) formatToPrecision(value float64) string {
	multiplier := math.Pow10(chart.Precision)
	multiplied := math.Round(value * multiplier)
	roundedValue := multiplied / multiplier

	result := fmt.Sprintf("%v", roundedValue)
	expectedLength := len(fmt.Sprintf("%.0f", multiplied))
	if chart.Precision > 0 {
		expectedLength++
	}

	for len(result) < expectedLength {
		if !strings.Contains(result, ".") {
			result += ".0"
		} else {
			result += "0"
		}
	}

	return result
}
