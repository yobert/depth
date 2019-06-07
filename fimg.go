package main

import (
	"fmt"
	"image"
)

type Fimg struct {
	Pix  []float64
	Rect image.Rectangle
}

type Fcolor struct {
	R, G, B, A float64
}

func NewFimg(r image.Rectangle) *Fimg {
	pix := make([]float64, r.Dx()*r.Dy()*4)
	return &Fimg{
		Pix:  pix,
		Rect: r,
	}
}

func (f Fimg) PixOffset(x, y int) int {
	w := f.Rect.Dx()
	y -= f.Rect.Min.Y
	x -= f.Rect.Min.X
	return y*w*4 + x*4
}
func (f Fimg) Set(x, y int, c Fcolor) {
	if !(image.Point{x, y}.In(f.Rect)) {
		return
	}
	i := f.PixOffset(x, y)
	s := f.Pix[i : i+4 : i+4]
	s[0] = c.R
	s[1] = c.G
	s[2] = c.B
	s[3] = c.A
}
func (f Fimg) Get(x, y int) Fcolor {
	if !(image.Point{x, y}.In(f.Rect)) {
		return Fcolor{}
	}
	i := f.PixOffset(x, y)
	s := f.Pix[i : i+4 : i+4]
	return Fcolor{
		s[0],
		s[1],
		s[2],
		s[3],
	}
}

func (f Fimg) Blot(x, y float64, c Fcolor) {
	w, h := f.Rect.Dx(), f.Rect.Dy()

	x += 1
	y += 1

	x *= 0.5
	y *= 0.5

	x *= float64(w)
	y *= float64(h)

	ix := f.Rect.Min.X + int(x)
	iy := f.Rect.Min.Y + int(y)

	oldc := f.Get(ix, iy)

	newc := Fcolor{
		//		oldc.R + c.R,
		//		oldc.G + c.G,
		//		oldc.B + c.B,
		//		oldc.A + c.A,
		c.R*c.A + oldc.R,
		c.G*c.A + oldc.G,
		c.B*c.A + oldc.B,
		0,
	}

	f.Set(ix, iy, newc)
}

func (f Fimg) ToNRGBA() *image.NRGBA {
	r := f.Rect

	r.Max.X -= r.Min.X
	r.Min.X = 0

	r.Max.Y -= r.Min.Y
	r.Min.Y = 0

	img := image.NewNRGBA(r)

	sum := 0.0

	for i, v := range f.Pix {
		if i%4 == 0 {
			sum += v
		}
		c := int(v * 255.0)
		if c < 0 {
			c = 0
		}
		if c > 255 {
			c = 255
		}
		img.Pix[i] = uint8(c)
	}
	//fmt.Println(sum)
	_ = fmt.Println
	return img
}

func (f Fimg) Bg() {
	for i := 0; i < len(f.Pix); i += 4 {
		/*		a := f.Pix[i+3]
				if a > 1 {
					f.Pix[i+0] /= a
					f.Pix[i+1] /= a
					f.Pix[i+2] /= a
				} else {
					f.Pix[i+0] *= a
					f.Pix[i+1] *= a
					f.Pix[i+2] *= a
				}*/
		f.Pix[i+3] = 1
	}
}
