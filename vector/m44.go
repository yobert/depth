package vector

import "fmt"

type M44 [16]float64

func IdentityM44() M44 {
	return M44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func TranslateM44(v V3) M44 {
	return M44{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v.X, v.Y, v.Z, 1}
}

func ScaleM44(v V3) M44 {
	return M44{
		v.X, 0, 0, 0,
		0, v.Y, 0, 0,
		0, 0, v.Z, 0,
		0, 0, 0, 1}
}

func (a M44) Mult(b M44) M44 {
	return M44{
		a[0]*b[0] + a[4]*b[1] + a[8]*b[2] + a[12]*b[3],
		a[1]*b[0] + a[5]*b[1] + a[9]*b[2] + a[13]*b[3],
		a[2]*b[0] + a[6]*b[1] + a[10]*b[2] + a[14]*b[3],
		a[3]*b[0] + a[7]*b[1] + a[11]*b[2] + a[15]*b[3],

		a[0]*b[4] + a[4]*b[5] + a[8]*b[6] + a[12]*b[7],
		a[1]*b[4] + a[5]*b[5] + a[9]*b[6] + a[13]*b[7],
		a[2]*b[4] + a[6]*b[5] + a[10]*b[6] + a[14]*b[7],
		a[3]*b[4] + a[7]*b[5] + a[11]*b[6] + a[15]*b[7],

		a[0]*b[8] + a[4]*b[9] + a[8]*b[10] + a[12]*b[11],
		a[1]*b[8] + a[5]*b[9] + a[9]*b[10] + a[13]*b[11],
		a[2]*b[8] + a[6]*b[9] + a[10]*b[10] + a[14]*b[11],
		a[3]*b[8] + a[7]*b[9] + a[11]*b[10] + a[15]*b[11],

		a[0]*b[12] + a[4]*b[13] + a[8]*b[14] + a[12]*b[15],
		a[1]*b[12] + a[5]*b[13] + a[9]*b[14] + a[13]*b[15],
		a[2]*b[12] + a[6]*b[13] + a[10]*b[14] + a[14]*b[15],
		a[3]*b[12] + a[7]*b[13] + a[11]*b[14] + a[15]*b[15]}
}

// OpenGL style matrix multiplication:
// it is actually Mult() with the operands switched
func (a M44) MultX(b M44) (o M44) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			o[i*4+j] =
				a[i*4+0]*b[0*4+j] +
					a[i*4+1]*b[1*4+j] +
					a[i*4+2]*b[2*4+j] +
					a[i*4+3]*b[3*4+j]
		}
	}
	return
}

func (m M44) MultV3(vec V3) (out V3) {
	out.X = vec.X*m[0] + vec.Y*m[4] + vec.Z*m[8] + m[12]
	out.Y = vec.X*m[1] + vec.Y*m[5] + vec.Z*m[9] + m[13]
	out.Z = vec.X*m[2] + vec.Y*m[6] + vec.Z*m[10] + m[14]
	return
}

func (m M44) MultV4(vec V4) (out V4) {
	out.X = vec.X*m[0] + vec.Y*m[4] + vec.Z*m[8] + vec.W*m[12]
	out.Y = vec.X*m[1] + vec.Y*m[5] + vec.Z*m[9] + vec.W*m[13]
	out.Z = vec.X*m[2] + vec.Y*m[6] + vec.Z*m[10] + vec.W*m[14]
	out.W = vec.X*m[3] + vec.Y*m[7] + vec.Z*m[11] + vec.W*m[15]
	return
}

func (m M44) Inverse() M44 {
	a0 := m[0]*m[5] - m[4]*m[1]
	a1 := m[0]*m[9] - m[8]*m[1]
	a2 := m[0]*m[13] - m[12]*m[1]
	a3 := m[4]*m[9] - m[8]*m[5]
	a4 := m[4]*m[13] - m[12]*m[5]
	a5 := m[8]*m[13] - m[12]*m[9]
	b0 := m[2]*m[7] - m[6]*m[3]
	b1 := m[2]*m[11] - m[10]*m[3]
	b2 := m[2]*m[15] - m[14]*m[3]
	b3 := m[6]*m[11] - m[10]*m[7]
	b4 := m[6]*m[15] - m[14]*m[7]
	b5 := m[10]*m[15] - m[14]*m[11]

	det := a0*b5 - a1*b4 + a2*b3 + a3*b2 - a4*b1 + a5*b0

	if det == 0.0 {
		return IdentityM44()
	}

	id := 1.0 / det

	return M44{
		id * (+m[5]*b5 - m[9]*b4 + m[13]*b3),
		id * (-m[1]*b5 + m[9]*b2 - m[13]*b1),
		id * (+m[1]*b4 - m[5]*b2 + m[13]*b0),
		id * (-m[1]*b3 + m[5]*b1 - m[9]*b0),
		id * (-m[4]*b5 + m[8]*b4 - m[12]*b3),
		id * (+m[0]*b5 - m[8]*b2 + m[12]*b1),
		id * (-m[0]*b4 + m[4]*b2 - m[12]*b0),
		id * (+m[0]*b3 - m[4]*b1 + m[8]*b0),
		id * (+m[7]*a5 - m[11]*a4 + m[15]*a3),
		id * (-m[3]*a5 + m[11]*a2 - m[15]*a1),
		id * (+m[3]*a4 - m[7]*a2 + m[15]*a0),
		id * (-m[3]*a3 + m[7]*a1 - m[11]*a0),
		id * (-m[6]*a5 + m[10]*a4 - m[14]*a3),
		id * (+m[2]*a5 - m[10]*a2 + m[14]*a1),
		id * (-m[2]*a4 + m[6]*a2 - m[14]*a0),
		id * (+m[2]*a3 - m[6]*a1 + m[10]*a0),
	}
}

// M33 will truncate the 4x4 matrix down to a 3x3
func (m M44) M33() M33 {
	return M33{
		m[0], m[1], m[2],
		m[4], m[5], m[6],
		m[8], m[9], m[10],
	}
}

func (a M44) TranslatePart() V3 {
	return V3{a[12], a[13], a[14]}
}

func (a M44) ScalePart() V3 {
	return V3{a[0], a[5], a[10]}
}

func (m M44) String() string {
	return fmt.Sprintf("[\t%.2f\t%.2f\t%.2f\t%.2f\n\t%.2f\t%.2f\t%.2f\t%.2f\n\t%.2f\t%.2f\t%.2f\t%.2f\n\t%.2f\t%.2f\t%.2f\t%.2f\t]\n",
		m[0], m[4], m[8], m[12],
		m[1], m[5], m[9], m[13],
		m[2], m[6], m[10], m[14],
		m[3], m[7], m[11], m[15])
}
