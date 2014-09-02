package system

type Platform struct {
	Cores []*Core
}

type Core struct {
	Time  []float64
	Power []float64
}

func NewPlatform(tables ...Table) (*Platform, error) {
	platform := &Platform{
	}

	return platform, nil
}
