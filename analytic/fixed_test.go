package analytic

import (
	"fmt"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
)

func TestFixedNew(t *testing.T) {
	const (
		nc = 2
	)

	temperature, _ := loadFixed(nc)

	assert.Equal(temperature.nc, uint(nc), t)
	assert.Equal(temperature.nn, uint(4*nc+12), t)

	assert.Close(temperature.D, fixtureD, 1e-14, t)

	assert.Close(temperature.E, fixtureE, 1e-9, t)
	assert.Close(temperature.F, fixtureF, 1e-9, t)
}

func TestFixedCompute(t *testing.T) {
	const (
		nc = 2
	)

	temperature, P := loadFixed(nc)
	Q := temperature.Compute(P)

	assert.Close(Q, fixtureQ, 1e-12, t)
}

func TestFixedComputeWithStatic(t *testing.T) {
	const (
		nc = 2
	)

	temperature, P := loadFixed(nc)
	noop := func([]float64, []float64) {}
	Q := temperature.ComputeWithStatic(P, noop)

	assert.Close(Q, fixtureQ, 1e-12, t)
}

func BenchmarkFixedCompute002(b *testing.B) {
	const (
		nc = 2
	)

	temperature, P := loadFixed(nc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P)
	}
}

func BenchmarkFixedComputeWithStatic002(b *testing.B) {
	const (
		nc = 2
	)

	temperature, P := loadFixed(nc)
	noop := func([]float64, []float64) {}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.ComputeWithStatic(P, noop)
	}
}

func BenchmarkFixedCompute032(b *testing.B) {
	const (
		nc = 32
		ns = 1000
	)

	temperature, _ := loadFixed(nc)
	P := random(nc*ns, 0, 20)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P)
	}
}

func loadFixed(nc uint) (*Fixed, []float64) {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%03d.json", nc)), config)
	temperature, _ := NewFixed(config)
	return temperature, append([]float64(nil), fixtureP...)
}
