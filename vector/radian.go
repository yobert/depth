package vector

import "math"

type Radian float64

func Sin(r Radian) float64 { return math.Sin((float64)(r)) }
func Cos(r Radian) float64 { return math.Cos((float64)(r)) }
func Tan(r Radian) float64 { return math.Tan((float64)(r)) }

//func Asin2(x, y float64) Radian { return (Radian)(math.Asin2(x, y)) }
//func Acos2(x, y float64) Radian { return (Radian)(math.Acos2(x, y)) }
func Atan2(x, y float64) Radian { return (Radian)(math.Atan2(x, y)) }

func Asin(x float64) Radian { return (Radian)(math.Asin(x)) }
func Acos(x float64) Radian { return (Radian)(math.Acos(x)) }
func Atan(x float64) Radian { return (Radian)(math.Atan(x)) }

// Degree converts a radian into degrees
func (φ Radian) Degree() Degree {
	return (Degree)(φ / τ * 360)
}
