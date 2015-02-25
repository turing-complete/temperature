package temperature

import (
	"errors"

	"github.com/ready-steady/hotspot"
)

// Config represents the configuration of a particular problem.
type Config struct {
	// The configuration of the HotSpot model.
	hotspot.Config

	// The sampling interval of temperature analysis.
	TimeStep float64 // in seconds

	// The temperature of the ambience.
	Ambience float64 // in Kelvin
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.TimeStep <= 0 {
		return errors.New("the time step should be positive")
	}

	return nil
}
