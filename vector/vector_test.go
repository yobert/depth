package vector

import (
	"math"
	"testing"
)

var _precision = 0.00001

func feq(a, b float64) bool {
	delta := a - b
	if math.Abs(delta) < _precision {
		return true
	}
	return false
}
func fne(a, b float64) bool {
	return !feq(a, b)
}
func v3eq(a, b V3) bool {
	if feq(a.X, b.X) && feq(a.Y, b.Y) && feq(a.Z, b.Z) {
		return true
	}
	return false
}
func m33eq(a, b M33) bool {
	for i := 0; i < 9; i++ {
		if fne(a[i], b[i]) {
			return false
		}
	}
	return true
}
func qeq(a, b Q) bool {
	if feq(a.R, b.R) && feq(a.I, b.I) && feq(a.J, b.J) && feq(a.K, b.K) {
		return true
	}
	return false
}

func TestDegree(t *testing.T) {
	_precision = 0.0001

	var d Degree = 45
	var a Radian = 2 * math.Pi

	if fne((float64)(d.Radian()), math.Pi/4) {
		t.Error("Degree Radian()")
	}

	if fne((float64)(a.Degree()), 360) {
		t.Error("Radian Degree()")
	}
}

func TestEuler(t *testing.T) {
	e := Euler{τ / 4, 0, 0}
	v := V3{0, 1, 0}

	if !v3eq(e.M33().MultV3(v), V3{0, 0, 1}) {
		t.Error("Euler -> M33 -> Rotate")
	}

	if !v3eq(e.Q().Euler().Q().M33().MultV3(v), V3{0, 0, 1}) {
		t.Error("Euler -> Q -> M33 -> Rotate ")
	}
}

func TestV3(t *testing.T) {
	_precision = 0.00001

	v := V3{1, 1, 0}

	var (
		m44 M44
		m33 M33
	)

	if fne(V3{10, 20, 30}.Len(), 37.416573868) {
		t.Error("V3 length")
	}
	if fne(V3{1, 2, 3}.Dot(V3{4, -5, 6}), 12) {
		t.Error("V3 dot product")
	}
	if !v3eq(V3{3, -3, 1}.Cross(V3{4, 9, 2}), V3{-15, -2, 39}) {
		t.Error("V3 cross product")
	}

	m33 = RotateAxisM33(V3{0, 0, 1}, -math.Pi/2)
	if !v3eq(m33.MultV3(v), V3{1, -1, 0}) {
		t.Error("M33 MultV3()")
	}
	if !v3eq(m33.Transpose().MultV3(v), V3{-1, 1, 0}) {
		t.Error("M33 Transpose() MultV3()")
	}

	m44 = m33.M44()
	if !v3eq(m44.MultV3(v), V3{1, -1, 0}) {
		t.Error("M44 RotateAxis + Multiply")
	}
	if !v3eq(m44.Inverse().MultV3(v), V3{-1, 1, 0}) {
		t.Error("M44 Inverse")
	}
}

func TestQ(t *testing.T) {
	_precision = 0.0001

	v := V3{1, 1, 0}

	var (
		m33 M33
		q   Q
	)

	q = AxisAngleQ(V3{0, 0, 1}, -τ/4)
	m33 = q.M33()
	if !v3eq(m33.MultV3(v), V3{1, -1, 0}) {
		t.Error("Q AxisAngleQ() M33()")
	}

	q = IdentityQ()
	m33 = q.M33()
	if !v3eq(m33.MultV3(v), V3{1, 1, 0}) {
		t.Error("IdentityQ()")
	}

	q = IdentityQ()
	q = q.Mult(AxisAngleQ(V3{0, 0, 1}, -τ/4))
	m33 = q.M33()
	if !v3eq(m33.MultV3(v), V3{1, -1, 0}) {
		t.Error("Q IdentityQ() Mult() AxisAngleQ() M33()")
	}

	q = IdentityQ()
	q = q.Mult(AxisAngleQ(V3{0, 0, 1}, -τ/4))
	q = q.Mult(AxisAngleQ(V3{0, 0, 1}, -τ/4))
	m33 = q.M33()
	if !v3eq(m33.MultV3(v), V3{-1, -1, 0}) {
		t.Error("Q IdentityQ() Mult() AxisAngleQ() M33() x 2")
	}

	q = Q{5, 6, 7, 8}.Add(Q{8, 7, 6, 5})
	if !qeq(q, Q{13, 13, 13, 13}) {
		t.Error("Q Add()")
	}
	q = Q{5, 6, 7, 8}.Sub(Q{8, 7, 6, 5})
	if !qeq(q, Q{-3, -1, 1, 3}) {
		t.Error("Q Sub()")
	}
	q = Q{2, 3, 2, 3}.Mult(Q{3, 2, 3, 2})
	if !qeq(q, Q{-12, 8, 12, 18}) {
		t.Error("Q Mult()")
	}

	/*	q = IdentityQ()
		q = q.Add(V3{0, 0, τ / 4}.Q().Mult(q)) // -90 degrees on the +z axis
		m33 = q.M33()
		if !v3eq(m33.MultV3(v), V3{1, -1, 0}) {
			t.Error("V3.Q() Q.Add()")
		}*/

}

func TestRadian(t *testing.T) {
	r := Radian(π)
	d := Degree(180)
	if fne(float64(r.Degree()), float64(d)) {
		t.Error("Radian->Degree")
	}
	if fne(float64(d.Radian()), float64(r)) {
		t.Error("Degree->Radian")
	}
}
