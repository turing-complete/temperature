package power

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/support/assert"
)

func TestCompute(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")

	assert.Equal(power.Compute(schedule, Δt, 440), fixtureP, t)
	assert.Equal(power.Compute(schedule, Δt, 42), fixtureP[:2*42], t)
}

func TestProcess(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")
	cc, sc := uint(2), uint(schedule.Span/Δt)

	process := power.Process(schedule)

	P := make([]float64, sc*cc)
	for i := uint(0); i < sc; i++ {
		process(Δt*(0.5+float64(i)), P[i*cc:(i+1)*cc])
	}

	mismatches := 0
	for i := range P {
		if P[i] != fixtureP[i] {
			mismatches++
		}
	}
	assert.Equal(mismatches, 17, t)
}

func BenchmarkCompute(b *testing.B) {
	const (
		Δt = 1e-5
	)

	power, schedule := prepare("002_040")
	sc := uint(schedule.Span / Δt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		power.Compute(schedule, Δt, sc)
	}
}

func prepare(name string) (*Power, *time.Schedule) {
	platform, application, _ := system.Load(findFixture(fmt.Sprintf("%s.tgff", name)))

	power := New(platform, application)

	profile := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	schedule := list.Compute(profile.Mobility)

	return power, schedule
}

func findFixture(name string) string {
	const (
		fixturePath = "fixtures"
	)
	return path.Join(fixturePath, name)
}
