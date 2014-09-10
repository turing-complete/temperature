package power

import (
	"path"
	"testing"

	"github.com/go-eslab/persim/system"
	"github.com/go-eslab/persim/time"
	"github.com/go-math/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestCompute(t *testing.T) {
	plat, app, _ := system.Load(findFixture("002_040.tgff"))

	prof := system.NewProfile(plat, app)
	list := time.NewList(plat, app)
	sched := list.Compute(prof.Mobility)

	power := New(plat, app, 1e-3)

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
