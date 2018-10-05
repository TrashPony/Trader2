package Analyze

func PercentageCalculator(a, b float64) (result float64) {
	result = 100 - (a * 100 / b)
	return result
}

