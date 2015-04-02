package power

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
)

func TestSample(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")

	assert.Equal(power.Sample(schedule, Δt, 440), fixtureP, t)
	assert.Equal(power.Sample(schedule, Δt, 42), fixtureP[:2*42], t)
}

func TestProcess(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")
	nc, ns := uint(2), uint(schedule.Span/Δt)

	process := power.Process(schedule)

	P := make([]float64, ns*nc)
	for i := uint(0); i < ns; i++ {
		process(Δt*(0.5+float64(i)), P[i*nc:(i+1)*nc])
	}

	mismatches := 0
	for i := range P {
		if P[i] != fixtureP[i] {
			mismatches++
		}
	}
	assert.Equal(mismatches, 17, t)
}

func BenchmarkSample(b *testing.B) {
	const (
		Δt = 1e-5
	)

	power, schedule := prepare("002_040")
	ns := uint(schedule.Span / Δt)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		power.Sample(schedule, Δt, ns)
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
