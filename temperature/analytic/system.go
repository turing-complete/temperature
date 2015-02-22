package analytic

type system struct {
	// D = C**(-½)
	D []float64

	// A = U * diag(Λ) * U**T
	U []float64
	Λ []float64

	// E = exp(A * Δt) = U * diag(exp(λi * Δt)) * U**T
	E []float64

	// F = A**(-1) * (exp(A * Δt) - I) * B
	//   = U * diag((exp(λi * Δt) - 1) / λi) * U**T * B
	F []float64

	Qamb float64
}
