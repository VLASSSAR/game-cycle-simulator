package generator

import (
	"math/rand"
)

type ArrivalGenerator struct {
	rng *rand.Rand
}

func NewArrivalGenerator(rng *rand.Rand) *ArrivalGenerator {
	return &ArrivalGenerator{
		rng: rng,
	}
}

// GenerateArrivals генерирует моменты поступления ставок
// для пуассоновского процесса с интенсивностью lambda.
// Возвращается возрастающий список времён прихода ставок.
func (g *ArrivalGenerator) GenerateArrivals(lambda float64, n int) []float64 {
	arrivals := make([]float64, 0, n)

	currentTime := 0.0
	for i := 0; i < n; i++ {
		interArrival := g.rng.ExpFloat64() / lambda
		currentTime += interArrival
		arrivals = append(arrivals, currentTime)
	}

	return arrivals
}
