package analytic

import (
	"github.com/simulated-reality/hotspot"
)

// Config represents the configuration of a particular problem.
type Config struct {
	// The configuration of the thermal RC model.
	hotspot.Config

	// The ambient temperature.
	Ambience float64 // in Celsius or Kelvin

	// The sampling interval. The parameter is specific to the Fixed integrator.
	TimeStep float64 // in seconds
}
