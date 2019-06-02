package vector

type Line struct {
	Start V3
	End   V3
}

func (l Line) Lerp(v float64) V3 {
	return l.End.Sub(l.Start).Scale(v).Add(l.Start)
}
