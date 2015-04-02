// Package analytic provides an exponential integrator for systems of ordinary
// differential equations modeling temperature of electronic systems.
//
// The initial system is
//
//         dQ'
//     C * -- + G * (Q - Qamb) = M * P
//         dt
//
//     Q = M**T * Q'
//
// where C is a diagonal matrix of the thermal capacitance; G is a symmetric
// matrix of the thermal conductance; Q' is a vector of the temperature of the
// thermal nodes; Q is a subset of Q corresponding to the processing elements;
// Qamb is a vector of the ambient temperature; P is a vector of the power
// consumption of the processing elements; and M is a diagonal matrix whose
// diagonal elements equal to unity that maps the power dissipation of the
// processing elements onto the thermal nodes.
//
// The transformed system is
//
//     dS
//     -- = A * S + B * P
//     dt
//
//     Q = B**T * S + Qamb
//
// where
//
//     S = D**(-1) * (Q' - Qamb),
//     A = -D * G * D,
//     B = D * M, and
//     D = C**(-½).
//
// The eigendecomposition of A, which is real and symmetric, is
//
//     A = U * diag(Λ) * U**T.
//
// Assuming that the power dissipation from time 0 to time t is constant, the
// solution to the system at time t is as follows:
//
//     S(t) = E(t) * S(0) + F(t) * P(0)
//
// where
//
//     E(t) = exp(A * t) = U * diag(exp(λi * t)) * U**T and
//     F(t) = A**(-1) * (exp(A * t) - I) * B
//          = U * diag((exp(λi * t) - 1) / λi) * U**T * B.
package analytic
