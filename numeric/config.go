package numeric

import (
	"github.com/simulated-reality/hotspot"
)

// Config represents the configuration of a particular problem.
type Config struct {
	// The configuration of the HotSpot model.
	hotspot.Config

	// The temperature of the ambience.
	Ambience float64 // in Celsius or Kelvin
}
