package analytic

import (
	"github.com/turing-complete/hotspot"
)

// Config is a configuration of temperature analysis.
type Config struct {
	// The thermal RC model.
	hotspot.Config

	// The ambient temperature.
	Ambience float64 // in Kelvin

	// The sampling interval. The parameter is specific to the Fixed integrator.
	TimeStep float64 // in seconds
}
