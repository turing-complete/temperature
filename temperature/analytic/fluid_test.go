package analytic

import (
	"fmt"
	"math"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/generator"
	"github.com/ready-steady/probability/uniform"
)

func TestFluidNew(t *testing.T) {
	const (
		nc = 2
	)

	temperature, _ := loadFluid(nc)

	assert.Equal(temperature.nc, uint(nc), t)
	assert.Equal(temperature.nn, uint(4*nc+12), t)

	assert.EqualWithin(temperature.D, fixtureD, 1e-14, t)

	assert.EqualWithin(abs(temperature.U), abs(fixtureU), 1e-9, t)
	assert.EqualWithin(temperature.Λ, fixtureΛ, 1e-9, t)
}

func TestFluidCompute(t *testing.T) {
	const (
		nc = 2
	)

	temperature, config := loadFluid(nc)
	ns := uint(len(fixtureP) / nc)

	time := make([]float64, ns)
	for i := range time {
		time[i] = config.TimeStep
	}

	Q := temperature.Compute(fixtureP, time)

	assert.EqualWithin(Q, fixtureQ, 1e-12, t)
}

func BenchmarkFluidCompute002(b *testing.B) {
	const (
		nc = 2
	)

	temperature, config := loadFluid(nc)
	ns := uint(len(fixtureP) / nc)

	time := make([]float64, ns)
	for i := range time {
		time[i] = config.TimeStep
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(fixtureP, time)
	}
}

func BenchmarkFluidCompute032(b *testing.B) {
	const (
		nc = 32
		ns = 1000
	)

	temperature, config := loadFluid(nc)
	P := probability.Sample(uniform.New(0, 1), generator.New(0), nc*ns)

	time := make([]float64, ns)
	for i := range time {
		time[i] = config.TimeStep
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P, time)
	}
}

func abs(A []float64) []float64 {
	B := make([]float64, len(A))

	for i := range B {
		B[i] = math.Abs(A[i])
	}

	return B
}

func loadFluid(nc uint) (*Fluid, *Config) {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%03d.json", nc)), config)
	temperature, _ := NewFluid(config)
	return temperature, config
}
