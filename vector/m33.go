package vector

import (
	"fmt"
	"math"
)

type M33 [9]float64

func IdentityM33() M33 {
	return M33{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

// RotateAxisM33 returns a matrix that will rotate about axis by angle degrees.
// if you were looking down the axis, the rotation would be counterclockwise.
func RotateAxisM33(axis V3, angle Radian) M33 {

	l2 := axis.Dot(axis)
	l1 := math.Sqrt(l2)

	c := Cos(angle)
	s := Sin(angle)

	return M33{
		(axis.X*axis.X + (axis.Y*axis.Y+axis.Z*axis.Z)*c) / l2,
		(axis.X*axis.Y*(1.0-c) + axis.Z*l1*s) / l2,
		(axis.X*axis.Z*(1.0-c) - axis.Y*l1*s) / l2,

		(axis.X*axis.Y*(1.0-c) - axis.Z*l1*s) / l2,
		(axis.Y*axis.Y + (axis.X*axis.X+axis.Z*axis.Z)*c) / l2,
		(axis.Y*axis.Z*(1.0-c) + axis.X*l1*s) / l2,

		(axis.X*axis.Z*(1.0-c) + axis.Y*l1*s) / l2,
		(axis.Y*axis.Z*(1.0-c) - axis.X*l1*s) / l2,
		(axis.Z*axis.Z + (axis.X*axis.X+axis.Y*axis.Y)*c) / l2,
	}
}

func (m M33) Transpose() M33 {
	return M33{
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8]}
}

func (a M33) Determinant() float64 {
	return a[0]*a[4]*a[8] + a[2]*a[3]*a[7] + a[1]*a[5]*a[6] -
		a[2]*a[4]*a[6] - a[1]*a[3]*a[8] - a[0]*a[5]*a[7]
}

func (a M33) Mult(b M33) M33 {
	return M33{
		a[0]*b[0] + a[3]*b[1] + a[6]*b[2],
		a[1]*b[0] + a[4]*b[1] + a[7]*b[2],
		a[2]*b[0] + a[5]*b[1] + a[8]*b[2],

		a[0]*b[3] + a[3]*b[4] + a[6]*b[5],
		a[1]*b[3] + a[4]*b[4] + a[7]*b[5],
		a[2]*b[3] + a[5]*b[4] + a[8]*b[5],

		a[0]*b[6] + a[3]*b[7] + a[6]*b[8],
		a[1]*b[6] + a[4]*b[7] + a[7]*b[8],
		a[2]*b[6] + a[5]*b[7] + a[8]*b[8]}
}

// OpenGL style matrix multiplication:
// it is actually Mult() with the operands switched
func (a M33) MultX(b M33) (o M33) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			o[i*3+j] =
				a[i*3+0]*b[0*3+j] +
					a[i*3+1]*b[1*3+j] +
					a[i*3+2]*b[2*3+j]
		}
	}
	return
}

func (m M33) MultV3(v V3) V3 {
	return V3{
		m[0]*v.X + m[3]*v.Y + m[6]*v.Z,
		m[1]*v.X + m[4]*v.Y + m[7]*v.Z,
		m[2]*v.X + m[5]*v.Y + m[8]*v.Z}
}

func (m M33) M44() M44 {
	return M44{
		m[0], m[1], m[2], 0,
		m[3], m[4], m[5], 0,
		m[6], m[7], m[8], 0,
		0, 0, 0, 1}
}

// Q returns a quaternion from a rotation matrix.  The matrix is assumed
// to be orthogonal and special in that it's determinant is 1.
// I think this means any "just rotation" matrix should work.
func (m M33) Q() Q {
	t := m[0] + m[4] + m[8]

	// we have to do some special casing to not divide by zero.
	// http://www.euclideanspace.com/maths/geometry/rotations/conversions/matrixToQuaternion/index.htm
	if t > 0 {
		s := math.Sqrt(t+1) * 2
		return Q{
			0.25 * s,
			(m[5] - m[7]) / s,
			(m[6] - m[2]) / s,
			(m[1] - m[3]) / s,
		}
	} else if m[0] > m[4] && m[0] > m[8] {
		s := math.Sqrt(1+m[0]-m[4]-m[8]) * 2
		return Q{
			(m[5] - m[7]) / s,
			0.25 * s,
			(m[3] + m[1]) / s,
			(m[6] + m[2]) / s,
		}
	} else if m[4] > m[8] {
		s := math.Sqrt(1+m[4]-m[0]-m[8]) * 2
		return Q{
			(m[6] - m[2]) / s,
			(m[3] + m[1]) / s,
			0.25 * s,
			(m[7] + m[5]) / s,
		}
	}

	s := math.Sqrt(1+m[8]-m[0]-m[4]) * 2
	return Q{
		(m[1] - m[3]) / s,
		(m[6] + m[2]) / s,
		(m[7] + m[5]) / s,
		0.25 * s,
	}
}

func (m M33) String() string {
	return fmt.Sprintf("[\t%.2f\t%.2f\t%.2f\n\t%.2f\t%.2f\t%.2f\n\t%.2f\t%.2f\t%.2f\t]",
		m[0], m[3], m[6],
		m[1], m[4], m[7],
		m[2], m[5], m[8])
}
