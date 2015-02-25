package numeric

import (
	"fmt"
	"path"
	"testing"

	"github.com/ready-steady/support/assert"
	"github.com/ready-steady/support/fixture"
)

func TestNew(t *testing.T) {
	temperature := load("002")

	assert.Equal(temperature.Cores, uint(2), t)
	assert.Equal(temperature.Nodes, uint(4*2+12), t)

	assert.Equal(temperature.system.A, fixtureA, t)
	assert.Equal(temperature.system.B, fixtureB, t)
}

func load(path string) *Temperature {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%s.json", path)), config)
	temperature, _ := New(config)
	return temperature
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
