package temperature

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/support/assert"
)

func TestComputeTransient(t *testing.T) {
	temperature, _ := Load(findFixture("002.json"))

	cc := uint32(2)
	sc := uint32(len(fixtureP)) / cc

	Q := make([]float64, cc*sc)
	temperature.ComputeTransient(fixtureP, Q, nil, sc)

	assert.AlmostEqual(Q, fixtureQ, t)
}

func BenchmarkComputeTransient(b *testing.B) {
	temperature, _ := Load(findFixture("032.json"))

	cc := uint32(32)
	sc := uint32(1000)
	nc := temperature.Nodes

	P := probability.Sample(uniform.New(0, 1), cc*sc)
	Q := make([]float64, cc*sc)
	S := make([]float64, nc*sc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.ComputeTransient(P, Q, S, sc)
	}
}
