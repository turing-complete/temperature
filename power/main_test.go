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

func TestPowerCompute(t *testing.T) {
	plat, app, _ := system.LoadTGFF(findFixture("002_040.tgff"))

	prof := system.NewProfile(plat, app)
	list := time.NewList(plat, app)
	sched := list.Compute(prof.Mobility)

	power := New(plat, app, 1e-3)
	P := power.Compute(sched)

	assert.Equal(P.Data, fixturePData, t)
}

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}
