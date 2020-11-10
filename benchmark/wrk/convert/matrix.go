package convert

import (
	"fmt"
	"strings"
)

// valueMatrix is a string matrix representation of the terminal output.
type valueMatrix [][]string

func newValueMatrix(output string) (matrix valueMatrix) {
	lines := splitAndFilterEmpty(output, "\n")
	for _, line := range lines {
		cells := splitAndFilterEmpty(line, " ")
		matrix = append(matrix, cells)
	}
	return
}

func (matrix valueMatrix) Get(rowIndex int, columnIndex int) (string, error) {
	if isValidIndex(rowIndex, len(matrix)) {
		row := matrix[rowIndex]
		if isValidIndex(columnIndex, len(row)) {
			return row[columnIndex], nil
		}
	}
	return "", fmt.Errorf("Invalid matrix index: (%d,%d)", rowIndex, columnIndex)
}

func isValidIndex(index int, length int) bool {
	return index < length && index >= 0
}

func splitAndFilterEmpty(str string, sep string) (items []string) {
	for _, item := range strings.Split(str, sep) {
		trimmedItem := strings.TrimSpace(item)
		if trimmedItem != "" {
			items = append(items, trimmedItem)
		}
	}
	return
}
