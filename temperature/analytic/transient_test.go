package analytic

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/support/assert"
)

func TestCompute(t *testing.T) {
	temperature := load("002")
	ns := uint(len(fixtureP)) / 2

	Q := temperature.Compute(fixtureP, ns)

	assert.EqualWithin(Q, fixtureQ, 1e-12, t)
}

func BenchmarkCompute002(b *testing.B) {
	temperature := load("002")
	ns := uint(len(fixtureP)) / 2

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(fixtureP, ns)
	}
}

func BenchmarkCompute032(b *testing.B) {
	temperature := load("032")
	nc, ns := uint(32), uint(1000)

	P := probability.Sample(uniform.New(0, 1), nc*ns)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P, ns)
	}
}
