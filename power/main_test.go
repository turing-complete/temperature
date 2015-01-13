package power

import (
	"path"
	"testing"

	"github.com/ready-steady/persim/system"
	"github.com/ready-steady/persim/time"
	"github.com/ready-steady/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestCompute(t *testing.T) {
	platform, application, _ := system.Load(findFixture("002_040.tgff"))

	prof := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	sched := list.Compute(prof.Mobility)

	power := New(platform, application, 1e-3)

	P := make([]float64, 2*440)
	power.Compute(sched, P, 440)
	assert.Equal(P, fixturePData, t)

	P = make([]float64, 2*42)
	power.Compute(sched, P, 42)
	assert.Equal(P, fixturePData[:2*42], t)
}

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}
