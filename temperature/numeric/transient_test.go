package numeric

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/support/assert"
)

func TestCompute(t *testing.T) {
	temperature := load("002")
	cc, sc, Δt := uint(2), uint(440), 1e-3

	power := smooth(fixtureP, cc, sc, Δt)
	time := time(Δt, sc)

	Q, err := temperature.Compute(power, time)

	assert.Success(err, t)
	assert.EqualWithin(Q, fixtureQ, 2e-10, t)
}

func BenchmarkCompute002(b *testing.B) {
	temperature := load("002")
	cc, sc, Δt := uint(2), uint(440), 1e-3

	power := smooth(fixtureP, cc, sc, Δt)
	time := time(Δt, sc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(power, time)
	}
}

func BenchmarkCompute032(b *testing.B) {
	temperature := load("032")
	cc, sc, Δt := uint(32), uint(1000), 1e-3

	P := probability.Sample(uniform.New(0, 20), cc*sc)
	power := smooth(P, cc, sc, Δt)
	time := time(Δt, sc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(power, time)
	}
}

func smooth(P []float64, cc, sc uint, Δt float64) func(float64, []float64) {
	return func(time float64, power []float64) {
		k := uint(time / Δt)
		for i := uint(0); i < cc; i++ {
			power[i] = P[k*cc+i]
		}
	}
}

func time(Δt float64, sc uint) []float64 {
	time := make([]float64, sc)
	for i := uint(0); i < sc; i++ {
		time[i] = float64(i) * Δt
	}
	return time
}
