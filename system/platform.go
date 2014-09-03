package system

type Platform struct {
	Cores []Core
}

type Core struct {
	ID    uint32
	Time  []float64
	Power []float64
}
