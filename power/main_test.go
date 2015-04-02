package power

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
)

func TestPartition(t *testing.T) {
	const (
		n = 10
		ε = 1e-14
	)

	power, schedule := prepare("002_040")

	points := make([]float64, n)
	for i := range points {
		points[i] = float64(i+1) * schedule.Span / n
	}
	P, ΔT, index := power.Partition(schedule, points, ε)

	assert.Equal(P, fixturePartition.P, t)
	assert.EqualWithin(ΔT, fixturePartition.ΔT, 1e-15, t)
	assert.Equal(index, fixturePartition.index, t)

	Σ := 0.0
	for i, j := uint(0), uint(0); i < n; i++ {
		for ; j < index[i]; j++ {
			Σ += ΔT[j]
		}
		assert.EqualWithin(Σ, points[i], 1e-15, t)
	}
}

func TestProgress(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")
	nc, ns := uint(2), uint(schedule.Span/Δt)

	progress := power.Progress(schedule)

	P := make([]float64, ns*nc)
	for i := uint(0); i < ns; i++ {
		progress(Δt*(0.5+float64(i)), P[i*nc:(i+1)*nc])
	}

	mismatches := 0
	for i := range P {
		if P[i] != fixtureSample.P[i] {
			mismatches++
		}
	}
	assert.Equal(mismatches, 17, t)
}

func TestSample(t *testing.T) {
	const (
		Δt = 1e-3
	)

	power, schedule := prepare("002_040")

	assert.Equal(power.Sample(schedule, Δt, 440), fixtureSample.P, t)
	assert.Equal(power.Sample(schedule, Δt, 42), fixtureSample.P[:2*42], t)
}

func TestTraverse(t *testing.T) {
	const (
		ε = 1e-14
	)

	test := func(points, Δ []float64, index []uint) {
		a, b := traverse(points, ε)
		assert.Equal(a, Δ, t)
		assert.Equal(b, index, t)
	}

	test(
		[]float64{0, 1, 2, 3, 4, 1, 2, 3, 4, 5},
		[]float64{1, 1, 1, 1, 1},
		[]uint{0, 1, 2, 3, 4, 1, 2, 3, 4, 5},
	)

	test(
		[]float64{0, 0, 0, 2, 4, 0, 6, 8},
		[]float64{2, 2, 2, 2},
		[]uint{0, 0, 0, 1, 2, 0, 3, 4},
	)

	test(
		[]float64{15, 10, 6, 3, 3, 1},
		[]float64{2, 3, 4, 5},
		[]uint{4, 3, 2, 1, 1, 0},
	)
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
