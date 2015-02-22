package temperature

import (
	"path"
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig(findFixture("002.json"))

	assert.Success(err, t)

	assert.Equal(config.Floorplan, findFixture("002.flp"), t)
	assert.Equal(config.HotSpot.Config, findFixture("hotspot.config"), t)
	assert.Equal(config.HotSpot.Params, "", t)
	assert.Equal(config.TimeStep, 1e-3, t)
	assert.Equal(config.AmbientTemp, 318.15, t)
}

func findFixture(name string) string {
	return path.Join("fixtures", name)
}
