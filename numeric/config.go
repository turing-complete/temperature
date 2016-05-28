package numeric

import (
	"github.com/turing-complete/hotspot"
)

// Config is the configuration of a problem.
type Config struct {
	// The configuration of the HotSpot model.
	hotspot.Config

	// The temperature of the ambience.
	Ambience float64 // in Celsius or Kelvin
}
