package time

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestListCompute(t *testing.T) {
	plat, app, _ := system.Load(findFixture("002_040"))
	prof := system.NewProfile(plat, app)

	list := NewList(plat, app)
	sched := list.Compute(prof.Mobility)

	mapping := []uint{
		0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0,
		1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0,
	}

	order := []uint{
		0, 1, 2, 4, 7, 8, 5, 11, 3, 9, 15, 19, 20, 10, 12, 24, 30, 17, 22, 28, 29,
		18, 21, 27, 34, 26, 13, 16, 23, 31, 33, 25, 38, 39, 37, 35, 14, 36, 6, 32,
	}

	start := []float64{
		0.0000, 0.0230, 0.0350, 0.0910, 0.0470, 0.0790, 0.4220, 0.0680,
		0.0680, 0.1030, 0.1530, 0.0900, 0.1720, 0.3150, 0.4040, 0.1170,
		0.3360, 0.2040, 0.2590, 0.1380, 0.1530, 0.2720, 0.2190, 0.3360,
		0.1850, 0.3640, 0.3040, 0.2920, 0.2340, 0.2590, 0.2040, 0.3510,
		0.4250, 0.3510, 0.3040, 0.4000, 0.4150, 0.3880, 0.3770, 0.3880,
	}

	duration := []float64{
		0.0230, 0.0120, 0.0120, 0.0120, 0.0210, 0.0110, 0.0180, 0.0230,
		0.0110, 0.0140, 0.0190, 0.0150, 0.0130, 0.0210, 0.0110, 0.0210,
		0.0150, 0.0150, 0.0130, 0.0150, 0.0210, 0.0200, 0.0150, 0.0100,
		0.0190, 0.0130, 0.0110, 0.0120, 0.0250, 0.0150, 0.0180, 0.0190,
		0.0140, 0.0130, 0.0260, 0.0220, 0.0100, 0.0160, 0.0110, 0.0120,
	}

	finish := make([]float64, len(start))
	span := 0.0
	for i := range start {
		finish[i] = start[i] + duration[i]
		if span < finish[i] {
			span = finish[i]
		}
	}

	assert.Equal(sched.Mapping, mapping, t)
	assert.Equal(sched.Order, order, t)
	assert.EqualWithin(sched.Start, start, 1e-15, t)
	assert.EqualWithin(sched.Finish, finish, 1e-15, t)
	assert.EqualWithin(sched.Span, span, 1e-15, t)
}

func TestListRecompute(t *testing.T) {
	plat, app, _ := system.Load(findFixture("002_040"))
	prof := system.NewProfile(plat, app)

	list := NewList(plat, app)
	sched := list.Compute(prof.Mobility)

	delay := []float64{
		0.0352, 0.0831, 0.0585, 0.0550, 0.0917, 0.0286, 0.0757, 0.0754,
		0.0380, 0.0568, 0.0076, 0.0054, 0.0531, 0.0779, 0.0934, 0.0130,
		0.0569, 0.0469, 0.0012, 0.0337, 0.0162, 0.0794, 0.0311, 0.0529,
		0.0166, 0.0602, 0.0263, 0.0654, 0.0689, 0.0748, 0.0451, 0.0084,
		0.0229, 0.0913, 0.0152, 0.0826, 0.0538, 0.0996, 0.0078, 0.0443,
	}

	start := []float64{
		0.0000, 0.0582, 0.1533, 0.4349, 0.2238, 0.3855, 1.7419, 0.3365,
		0.3365, 0.5019, 0.6554, 0.4251, 0.6820, 1.2139, 1.6966, 0.5727,
		1.3128, 0.7837, 0.9856, 0.6067, 0.6554, 0.9998, 0.8456, 1.3128,
		0.7481, 1.4890, 1.1766, 1.0992, 0.8917, 0.9856, 0.7837, 1.3847,
		1.8648, 1.3847, 1.1766, 1.6373, 1.8010, 1.5810, 1.5622, 1.5810,
	}

	finish := make([]float64, len(start))
	for i := range finish {
		finish[i] = start[i] + (sched.Finish[i] - sched.Start[i]) + delay[i]
	}

	sched = list.Recompute(sched, delay)

	assert.EqualWithin(sched.Start, start, 1e-15, t)
	assert.EqualWithin(sched.Finish, finish, 2e-15, t)
}

func TestListRecomputeDummy(t *testing.T) {
	plat, app, _ := system.Load(findFixture("002_040"))
	prof := system.NewProfile(plat, app)

	list := NewList(plat, app)
	sched1 := list.Compute(prof.Mobility)
	sched2 := list.Recompute(sched1, make([]float64, len(sched1.Start)))

	assert.Equal(sched2.Start, sched1.Start, t)
	assert.Equal(sched2.Finish, sched1.Finish, t)
}

func BenchmarkListCompute_002_040(b *testing.B) {
	benchmarkCompute("002_040", b)
}

func BenchmarkListRecompute_002_040(b *testing.B) {
	benchmarkRecompute("002_040", b)
}

func BenchmarkListCompute_032_640(b *testing.B) {
	benchmarkCompute("032_640", b)
}

func BenchmarkListRecompute_032_640(b *testing.B) {
	benchmarkRecompute("032_640", b)
}

func benchmarkCompute(name string, b *testing.B) {
	plat, app, _ := system.Load(findFixture(name))

	prof := system.NewProfile(plat, app)
	list := NewList(plat, app)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list.Compute(prof.Mobility)
	}
}

func benchmarkRecompute(name string, b *testing.B) {
	plat, app, _ := system.Load(findFixture(name))

	prof := system.NewProfile(plat, app)
	list := NewList(plat, app)
	sched := list.Compute(prof.Mobility)
	delay := make([]float64, len(sched.Start))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		list.Recompute(sched, delay)
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
