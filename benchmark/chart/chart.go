package chart

import "io"

// Chart can render an svg chart.
type Chart interface {
	Render(target io.Writer)
}
