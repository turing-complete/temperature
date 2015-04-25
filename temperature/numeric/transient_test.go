package numeric

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestCompute002Fixed(t *testing.T) {
	const (
		nc = 2
		ns = 440
		Δt = 1e-3
	)

	temperature := load(nc)
	power := smooth(fixtureP, nc, ns, Δt)
	time := sequence(ns, Δt)

	Q, _, _ := temperature.Compute(power, time)

	assert.EqualWithin(Q, fixtureQ, 2e-10, t)
}

func TestCompute002Adaptive(t *testing.T) {
	const (
		nc = 2
		ns = 440
		Δt = 1e-3
	)

	temperature := load(nc)
	power := smooth(fixtureP, nc, ns, Δt)
	Q, time, _ := temperature.Compute(power, []float64{0, ns * Δt})

	assert.EqualWithin(Q, fixtureQTime, 1e-10, t)
	assert.EqualWithin(time, fixtureTime, 1e-14, t)
}

func BenchmarkCompute002Adaptive(b *testing.B) { benchmarkComputeAdaptive(2, 1000, 1e-3, b) }
func BenchmarkCompute032Adaptive(b *testing.B) { benchmarkComputeAdaptive(32, 1000, 1e-3, b) }

func benchmarkComputeAdaptive(nc, ns uint, Δt float64, b *testing.B) {
	temperature := load(nc)
	power := smooth(random(nc*ns, 0, 20), nc, ns, Δt)
	time := []float64{0, float64(ns) * Δt}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(power, time)
	}
}

func BenchmarkCompute002Fixed(b *testing.B) { benchmarkComputeFixed(2, 1000, 1e-3, b) }
func BenchmarkCompute032Fixed(b *testing.B) { benchmarkComputeFixed(32, 1000, 1e-3, b) }

func benchmarkComputeFixed(nc, ns uint, Δt float64, b *testing.B) {
	temperature := load(nc)
	power := smooth(random(nc*ns, 0, 20), nc, ns, Δt)
	time := sequence(ns, Δt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(power, time)
	}
}

func smooth(P []float64, nc, ns uint, Δt float64) func(float64, []float64) {
	return func(time float64, power []float64) {
		k := uint(time / Δt)
		if k >= ns {
			k = ns - 1
		}
		for i := uint(0); i < nc; i++ {
			power[i] = P[k*nc+i]
		}
	}
}
