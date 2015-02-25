package temperature

import (
	"path"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig(findFixture("002.json"))

	assert.Success(err, t)

	assert.Equal(config.HotSpot.Floorplan, findFixture("002.flp"), t)
	assert.Equal(config.HotSpot.Configuration, findFixture("hotspot.config"), t)
	assert.Equal(config.HotSpot.Parameters, "", t)
	assert.Equal(config.TimeStep, 1e-3, t)
	assert.Equal(config.Ambience, 318.15, t)
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
