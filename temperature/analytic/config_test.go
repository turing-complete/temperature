package analytic

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestConfigLoad(t *testing.T) {
	config, err := loadConfig(findFixture("002.json"))

	assert.Success(err, t)

	assert.Equal(config.Floorplan, findFixture("002.flp"), t)
	assert.Equal(config.HotSpot.Config, findFixture("hotspot.config"), t)
	assert.Equal(config.HotSpot.Params, "", t)
	assert.Equal(config.TimeStep, 1e-3, t)
	assert.Equal(config.AmbientTemp, 318.15, t)
}
