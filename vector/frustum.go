package vector

// Frustom represents a truncated pyramid view frustum.
type Frustum struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
	Near   float64
	Far    float64
}

func PerspectiveFrustum(y_fov Degree, x_ratio, near, far float64) (f Frustum) {

	side := Tan(y_fov.Radian() / 2)

	if x_ratio < 1.0 {
		f.Right = near * side
		f.Top = f.Right * (1.0 / x_ratio)
	} else {
		f.Top = near * side
		f.Right = f.Top * x_ratio
	}

	f.Bottom = -f.Top
	f.Left = -f.Right
	f.Near = near
	f.Far = far

	return
}

// M44 converts the frustum into a 4x4 matrix suitible for doing perspective transformations
func (f Frustum) M44() (m M44) {
	t1 := f.Near * 2.0
	t2 := f.Right - f.Left
	t3 := f.Top - f.Bottom
	t4 := f.Far - f.Near

	m[0] = t1 / t2
	m[1] = 0.0
	m[2] = 0.0
	m[3] = 0.0

	m[4] = 0.0
	m[5] = t1 / t3
	m[6] = 0.0
	m[7] = 0.0

	m[8] = (f.Right + f.Left) / t2
	m[9] = (f.Top + f.Bottom) / t3
	m[10] = (f.Far + f.Near) / -t4
	m[11] = -1.0

	m[12] = 0.0
	m[13] = 0.0
	m[14] = (-t1 * f.Far) / t4
	m[15] = 0.0

	return
}
