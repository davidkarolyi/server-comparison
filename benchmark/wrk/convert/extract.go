package convert

import (
	"fmt"
	"strconv"
	"unicode"
)

// fieldExtractor holds functionality to extract values from an output matrix.
type fieldExtractor struct {
	matrix valueMatrix
}

func newFieldExtractor(matrix valueMatrix) *fieldExtractor {
	return &fieldExtractor{
		matrix: matrix,
	}
}

func (extractor *fieldExtractor) ReqsPerSec() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Requests/sec:")
	if err != nil {
		return 0, err
	}

	value, err := extractor.matrix.Get(rowIndex, 1)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(value, 64)
}

func (extractor *fieldExtractor) TransferedBytesPerSec() (int64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Transfer/sec:")
	if err != nil {
		return 0, err
	}

	quantity, err := extractor.matrix.Get(rowIndex, 1)
	if err != nil {
		return 0, err
	}

	value, unit := trimUnit(quantity)
	unitMultiplier := getMultiplierOfUnit(unit, targetUnitBase, scaleBinary)

	dataSize, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return int64(dataSize * unitMultiplier), nil
}

func (extractor *fieldExtractor) LatencyReqsInStdevPerc() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Latency")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 4, targetUnitBase, scaleMetric)
}

func (extractor *fieldExtractor) LatencyAvgMS() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Latency")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 1, targetUnitMilli, scaleMetric)
}

func (extractor *fieldExtractor) LatencyStdevMS() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Latency")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 2, targetUnitMilli, scaleMetric)
}

func (extractor *fieldExtractor) LatencyMaxMS() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Latency")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 3, targetUnitMilli, scaleMetric)
}

func (extractor *fieldExtractor) ReqPerSecAvg() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Req/Sec")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 1, targetUnitBase, scaleMetric)
}

func (extractor *fieldExtractor) ReqPerSecMax() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Req/Sec")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 3, targetUnitBase, scaleMetric)
}

func (extractor *fieldExtractor) ReqPerSecStdev() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Req/Sec")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 2, targetUnitBase, scaleMetric)
}

func (extractor *fieldExtractor) ReqPerSecReqsInStdevPerc() (float64, error) {
	rowIndex, err := extractor.findRowThatStartsWithString("Req/Sec")
	if err != nil {
		return 0, err
	}
	return extractor.extractFloatWithUnit(rowIndex, 4, targetUnitBase, scaleMetric)
}

func (extractor *fieldExtractor) findRowThatStartsWithString(str string) (int, error) {
	for index := range extractor.matrix {
		value, _ := extractor.matrix.Get(index, 0)
		if value == str {
			return index, nil
		}
	}
	return 0, fmt.Errorf("Couldn't find row in the matrix that starts with: '%s'", str)
}

func (extractor *fieldExtractor) extractFloatWithUnit(
	rowIndex int,
	columnIndex int,
	targetUnit string,
	scale float64,
) (float64, error) {
	quantity, err := extractor.matrix.Get(rowIndex, columnIndex)
	if err != nil {
		return 0, err
	}

	value, unit := trimUnit(quantity)
	unitMultiplier := getMultiplierOfUnit(unit, targetUnit, scale)

	float, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}

	return float * unitMultiplier, nil
}

func trimUnit(quantity string) (value string, unit string) {
	for index, char := range quantity {
		digitOrDot := unicode.IsDigit(char) || char == '.'
		if !digitOrDot {
			return quantity[:index], quantity[index:]
		}
	}
	return quantity, ""
}
