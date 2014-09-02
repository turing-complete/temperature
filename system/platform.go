package system

type Platform struct {
	Cores []Core
}

type Core struct {
	Time  []float64
	Power []float64
}
