package logic

import "math"

func NormalFloat64(x int64, miu int64, sigma int64) float64 {
	return 1 / (math.Sqrt(2*math.Pi) * float64(sigma)) * math.Pow(math.E, -math.Pow(float64(x-miu), 2)/(2*math.Pow(float64(sigma), 2)))
}
