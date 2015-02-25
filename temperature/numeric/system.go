package numeric

type system struct {
	// A = -C**(-1) * G
	A []float64

	// B = C**(-1) * M
	B []float64

	Δt   float64
	Qamb float64
}
