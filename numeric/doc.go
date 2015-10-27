// Package numeric provides a Runge–Kutta integrator for systems of ordinary
// differential equations modeling temperature of multiprocessor systems.
//
// The system is
//
//         dS
//     C * -- + G * S = M * P
//         dt
//
//     Q = M**T * S + Qamb
//
// where C is a diagonal matrix of the thermal capacitance; G is a symmetric
// matrix of the thermal conductance; S is a vector of the state of the thermal
// nodes; Q is a vector of the temperature of the processing elements; P is a
// vector of the power consumption of the processing elements; M is a diagonal
// matrix whose diagonal elements equal to unity that maps the power dissipation
// of the processing elements onto the thermal nodes; and Qamb is a vector of
// the ambient temperature.
//
// The transformed system is
//
//     dS
//     -- = A * S + B * P
//     dt
//
//     Q = M**T * S + Qamb
//
// where
//
//     A = -C**(-1) * G and
//     B = C**(-1) * M.
package numeric
