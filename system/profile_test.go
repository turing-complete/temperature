package system

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestNewProfile(t *testing.T) {
	plat, app, _ := LoadTGFF(findFixture("002_040"))

	prof := NewProfile(plat, app)

	mobility := []float64{
		0.0000, 0.0000, 0.0000, 0.0445, 0.0000, 0.0140, 0.1145, 0.0000,
		0.0095, 0.0000, 0.0145, 0.0140, 0.0140, 0.0280, 0.0740, 0.0000,
		0.0280, 0.0140, 0.0225, 0.0000, 0.0000, 0.0225, 0.0140, 0.0310,
		0.0000, 0.0380, 0.0260, 0.0225, 0.0140, 0.0140, 0.0000, 0.0280,
		0.1490, 0.0340, 0.0225, 0.0520, 0.0280, 0.0435, 0.0380, 0.0380,
	}

	assert.AlmostEqual(prof.Mobility, mobility, t)
}
