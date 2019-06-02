package vector

import (
	"fmt"
	"math"
)

// Q is a quaternion.
// r = real part
// i, j, k = imaginary vector parts
type Q struct {
	R, I, J, K float64
}

// IdentityQ returns a new quaternion that does not do any rotating.
func IdentityQ() Q {
	return Q{1.0, 0.0, 0.0, 0.0}
}

// AxisAngleQ returns a quaternion representing a rotation around an axis.
// For this to work, the axis must be normalized.
func AxisAngleQ(axis V3, φ Radian) Q {

	φ = φ / 2

	return Q{
		Cos(φ),
		Sin(φ) * axis.X,
		Sin(φ) * axis.Y,
		Sin(φ) * axis.Z}
}

// Normalize will ensure the quaternion represents only a rotation.
// Good to do once in a while if you've done lots of floating point math.
func (q Q) Normalize() Q {
	l := math.Sqrt(q.R*q.R + q.I*q.I + q.J*q.J + q.K*q.K)
	if l == 0.0 {
		return IdentityQ()
	}
	return q.Scale(1.0 / l)
}

func (a Q) Mult(b Q) Q {
	return Q{
		a.R*b.R - a.I*b.I - a.J*b.J - a.K*b.K,
		a.R*b.I + a.I*b.R + a.J*b.K - a.K*b.J,
		a.R*b.J - a.I*b.K + a.J*b.R + a.K*b.I,
		a.R*b.K + a.I*b.J - a.J*b.I + a.K*b.R}
}
func (a Q) Add(b Q) Q {
	return Q{a.R + b.R, a.I + b.I, a.J + b.J, a.K + b.K}
}
func (a Q) Sub(b Q) Q {
	return Q{a.R - b.R, a.I - b.I, a.J - b.J, a.K - b.K}
}

func (q Q) Scale(s float64) Q {
	return Q{q.R * s, q.I * s, q.J * s, q.K * s}
}

// Euler will try to return a set of 3 rotations about the x, y, and z axis.
// x: bank
// y: heading
// z: attitude
//r	i	j	k
//w	x	y	z
func (q Q) Euler() Euler {
	gimbal := q.I*q.J + q.K*q.R
	if gimbal > 0.499 {
		return Euler{
			0,
			2 * Atan2(q.I, q.R),
			math.Pi / 2}
	}

	if gimbal < -0.499 {
		return Euler{
			0,
			-2 * Atan2(q.I, q.R),
			-math.Pi / 2}
	}

	sqi := q.I * q.I
	sqj := q.J * q.J
	sqk := q.K * q.K

	return Euler{
		Atan2(2*q.I*q.R-2*q.J*q.K, 1-2*sqi-2*sqk),
		Atan2(2*q.J*q.R-2*q.I*q.K, 1-2*sqj-2*sqk),
		Asin(2 * gimbal)}
}

// M33 converts a quaternion to a 3x3 rotation matrix
// math notation:
// R I J K
//
// according to wikipedia:
//1 - 2JJ - 2KK	    2IJ - 2KR	    2IK + 2JR
//    2IJ + 2KR	1 - 2II - 2KK	    2JK - 2IR
//    2IK - 2JR	    2JK + 2IR	1 - 2II - 2JJ
// unfortunately this looks like a mistake...
// the below code uses 1 + for cells [0] and [4]...
// it passes the tests but hell if I know what I'm doing
func (q Q) M33() M33 {
	return M33{
		1.0 - (2.0*q.J*q.J + 2.0*q.K*q.K),
		//		1.0 - (2.0*q.J*q.J - 2.0*q.K*q.K), // different sign!?
		2.0*q.I*q.J + 2.0*q.K*q.R,
		2.0*q.I*q.K - 2.0*q.J*q.R,
		2.0*q.I*q.J - 2.0*q.K*q.R,
		1.0 - (2.0*q.I*q.I + 2.0*q.K*q.K),
		//		1.0 - (2.0*q.I*q.I - 2.0*q.K*q.K), // different sign!?!?!
		2.0*q.J*q.K + 2.0*q.I*q.R,
		2.0*q.I*q.K + 2.0*q.J*q.R,
		2.0*q.J*q.K - 2.0*q.I*q.R,
		1.0 - (2.0*q.I*q.I + 2.0*q.J*q.J)} // this last line is ok  !?
}

func (q Q) String() string {
	return fmt.Sprintf("\t{%.4f,\t%.4f,\t%.4f,\t%.4f}", q.R, q.I, q.J, q.K)
}
