package convert

import (
	"strconv"
	"unicode/utf8"
)

var UnitSymbol = map[string]float64{
	"TB":  1000000000000,
	"GB": 1000000000,
	"MB": 1000000,
	"KB": 1000,
}

var siFactors = map[string]float64{
	"":  1e0,
	"k": 1e3,
	"M": 1e6,
	"G": 1e9,
	"T": 1e12,
	"P": 1e15,
	"E": 1e18,
	"Z": 1e21,
	"Y": 1e24,
	"K": 1e3,
	"B": 1e9,
}

func StringToFloat(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f, nil
	}
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError {
		return 0, err
	}
	symbol := s[len(s)-size : len(s)]
	factor, ok := siFactors[symbol]
	if !ok {
		return 0, err
	}
	f, e := strconv.ParseFloat(s[:len(s)-len(symbol)], 64)
	if e != nil {
		return 0, err
	}
	return f * factor, nil
}

