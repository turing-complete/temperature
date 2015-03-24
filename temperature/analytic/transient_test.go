package analytic

import (
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
)

func TestCompute(t *testing.T) {
	const (
		nc = 2
		ns = 440
	)

	temperature := load(nc)
	Q := temperature.Compute(fixtureP, ns)

	assert.EqualWithin(Q, fixtureQ, 1e-12, t)
}

func BenchmarkCompute002(b *testing.B) {
	const (
		nc = 2
		ns = 440
	)

	temperature := load(nc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(fixtureP, ns)
	}
}

func BenchmarkCompute032(b *testing.B) {
	const (
		nc = 32
		ns = 1000
	)

	temperature := load(nc)
	P := probability.Sample(uniform.New(0, 1), nc*ns)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P, ns)
	}
}
