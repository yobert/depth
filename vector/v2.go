package vector

import (
	"fmt"
	"math"
)

type V2 struct {
	X, Y float64
}

func (v V2) Dist(a V2) float64 {
	return v.Sub(a).Len()
}
func (v V2) Len() float64 {
	return math.Sqrt(v.LenSq())
}
func (v V2) LenSq() float64 {
	return v.Dot(v)
}
func (v V2) Dot(a V2) float64 {
	return v.X*a.X + v.Y*a.Y
}
func (v V2) Cross(a V2) float64 {
	// not sure if this is right.
	// https://github.com/pgkelley4/line-segments-intersect/blob/master/js/line-segments-intersect.js
	return v.X*a.Y - v.Y*a.X
}
func (v V2) Sub(a V2) V2 {
	return V2{v.X - a.X, v.Y - a.Y}
}
func (v V2) Add(a V2) V2 {
	return V2{v.X + a.X, v.Y + a.Y}
}
func (v V2) Scale(s float64) V2 {
	return V2{v.X * s, v.Y * s}
}
func (v V2) Normalize() V2 {
	l := v.Len()
	if l == 0.0 {
		return V2{}
	}
	return v.Scale(1.0 / l)
}

func (v V2) String() string {
	return fmt.Sprintf("%.2f %.2f", v.X, v.Y)
}
