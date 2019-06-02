package vector

// V4 is a 4 component vector (x, y, z, and w usually)
type V4 struct {
	X, Y, Z, W float64
}

func (v V4) HomogeneousToCartesian() V3 {
	if v.W == 0.0 {
		return V3{}
	}
	return V3{
		v.X / v.W,
		v.Y / v.W,
		v.Z / v.W}
}
