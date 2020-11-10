package convert

import (
	"math"
)

const scaleBinary float64 = 1024
const scaleMetric float64 = 1000

const targetUnitMilli = "m"
const targetUnitBase = ""

var unitPrefixPowers = map[string]float64{
	"n": -3,
	"u": -2,
	"m": -1,
	"k": 1,
	"K": 1,
	"M": 2,
	"G": 3,
	"T": 4,
	"P": 5,
}

func getMultiplierOfUnit(unit string, targetUnit string, scale float64) float64 {
	if len(unit) == 0 {
		return 1
	}

	unitPrefix := string(unit[0])
	power := unitPrefixPowers[unitPrefix] - unitPrefixPowers[targetUnit]

	return math.Pow(scale, power)
}
