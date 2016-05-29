package analytic

import (
	"github.com/turing-complete/hotspot"
)

// Config is the configuration of a problem.
type Config struct {
	// The configuration of the thermal RC model.
	hotspot.Config

	// The ambient temperature.
	Ambience float64 // in Kelvin

	// The sampling interval. The parameter is specific to the Fixed integrator.
	TimeStep float64 // in seconds
}
