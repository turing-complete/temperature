package power

import (
	"path"
	"testing"

	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestCompute(t *testing.T) {
	platform, application, _ := system.Load(findFixture("002_040.tgff"))
	profile := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	schedule := list.Compute(profile.Mobility)
	power, _ := New(platform, application, 1e-3)

	assert.Equal(power.Compute(schedule, 440), fixturePData, t)
	assert.Equal(power.Compute(schedule, 42), fixturePData[:2*42], t)
}

func BenchmarkCompute(b *testing.B) {
	const (
		Δt = 1e-5
	)

	platform, application, _ := system.Load(findFixture("002_040.tgff"))
	profile := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	schedule := list.Compute(profile.Mobility)
	power, _ := New(platform, application, Δt)

	sc := uint(schedule.Span / Δt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		power.Compute(schedule, sc)
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}
