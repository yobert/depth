package vector

// Euler represents three amounts of rotation, about the X, Y, and Z axis.
type Euler struct {
	X, Y, Z Radian
}

// Q converts a Euler into a quaternion
func (e Euler) Q() Q {
	cx := Cos(e.X / 2)
	sx := Sin(e.X / 2)
	cy := Cos(e.Y / 2)
	sy := Sin(e.Y / 2)
	cz := Cos(e.Z / 2)
	sz := Sin(e.Z / 2)

	return Q{
		cx*cy*cz - sx*sy*sz,
		sx*cy*cz + cx*sy*sz,
		cx*sy*cz + sx*cy*sz,
		cx*cy*sz - sx*sy*cz}
}

// M33 converts a Euler into a 3x3 rotation matrix
func (e Euler) M33() M33 {
	return IdentityM33().
		Mult(RotateAxisM33(V3{0, 0, 1}, e.Z)).
		Mult(RotateAxisM33(V3{0, 1, 0}, e.Y)).
		Mult(RotateAxisM33(V3{1, 0, 0}, e.X))
}
