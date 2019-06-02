package vector

// M34 is a special bastard matrix type for holding a rotation (3x3) matrix along
// with a transpose (1x3) matrix in the same structure, with methods letting it
// behave like a 4x4 matrix by always assuming the bottom row is 0, 0, 0, 1.
type M34 [12]float64

func IdentityM34() M34 {
	return M34{
		1.0, 0.0, 0.0, 0.0,
		0.0, 1.0, 0.0, 0.0,
		0.0, 0.0, 1.0, 0.0}
}

func RotateTransposeM34(rotate M33, transpose V3) M34 {
	return M34{
		rotate[0], rotate[1], rotate[2], transpose.X,
		rotate[3], rotate[4], rotate[5], transpose.Y,
		rotate[6], rotate[7], rotate[8], transpose.Z}
}

/*func (m M443) MultV3(vec V3) {
	return V3{
		m[ 0]*vec.X + m[ 4]*vec.Y + m[ 8]*vec.Z + m[12],
		m[ 1]*vec.X + m[ 5]*vec.Y + m[ 9]*vec.Z + m[13],
		m[ 2]*vec.X + m[ 6]*vec.Y + m[10]*vec.Z + m[14]}
}

func (m M443) MultInverseV3(vec V3) {
	return V3{
		m[ 0]*vec.X + m[ 4]*vec.Y + m[ 8]*vec.Z + m[12],
		m[ 1]*vec.X + m[ 5]*vec.Y + m[ 9]*vec.Z + m[13],
		m[ 2]*vec.X + m[ 6]*vec.Y + m[10]*vec.Z + m[14]}*/
