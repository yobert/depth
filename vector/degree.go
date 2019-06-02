package vector

type Degree float64

// Radian() converts a degree into radians
func (φ Degree) Radian() Radian {
	return (Radian)(φ / 360 * τ)
}
