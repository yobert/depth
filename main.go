package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"

	"depth/vector"
)

func (f Fimg) BlotPoint(mv vector.M44, p vector.V3, c Fcolor) {
	p4 := vector.V4{p.X, p.Y, p.Z, 1}
	p4 = mv.MultV4(p4)
	pp := p4.HomogeneousToCartesian()
	if pp.Z > -1 && pp.Z < 1 {
		f.Blot(pp.X, pp.Y, c)
	}
}

func (img Fimg) BlotLine(mm, mv vector.M44, cam vector.Camera, ray vector.Line) {
	m := 0.3
	e := 3.0
	f := 1.3

	for i := 0; i < int(ray.Start.Dist(ray.End)*100000.0); i++ {
		v := ray.Lerp(rand.Float64())

		v = mm.MultV3(v)

		d := cam.Position.Dist(v)

		//a := 1.0 / (math.Pow(d, 2))
		//a *= 0.1
		a := 0.005

		r := m * math.Pow(math.Abs(f-d), e)
		w := v.Add(rndSphere(r))

		c := Fcolor{1, 1, 1, a}
		img.BlotPoint(mv, w, c)
	}
}

func rndSphere(r float64) vector.V3 {
	v := vector.V3{
		rand.Float64()*2.0 - 1.0,
		rand.Float64()*2.0 - 1.0,
		rand.Float64()*2.0 - 1.0,
	}
	v = v.Normalize()
	return v.Scale(r)
}

func main() {

	frame := 0


	for rot := vector.Radian(0); rot < 2 * math.Pi; rot += 0.001 {

		mm := vector.RotateAxisM33(vector.V3{0, -1, 0}, rot)

		rand.Seed(666)

		s := 500

		acc := NewFimg(image.Rect(-s, -s, s, s))

		cam := vector.Camera{
			Width:  float64(acc.Rect.Dx()),
			Height: float64(acc.Rect.Dy()),

			YFov: 60,
			Near: 0.1,
			Far:  100,

			Position: vector.V3{0, 0, 2.2},
			RotAxis:  vector.Euler{
				//			X: 0.01,
			},
		}

		cam.SetupViewProjection()
		cam.SetupModelView()

		mv := cam.ModelView.MultX(cam.Projection)

		if false {
			for x := -1.0; x < 1.0; x += 0.1 {
				for y := -1.0; y < 1.0; y += 0.1 {
					for z := -1.0; z < 1.0; z += 0.1 {
						c := Fcolor{0.5, 0.5, 0.5, 1}
						if x > 0.8 {
							c.R = 1
						}
						if y > 0.8 {
							c.G = 1
						}
						if z > 0.8 {
							c.B = 1
						}
						p := vector.V4{x, y, z, 1}
						p = mv.MultV4(p)
						pp := p.HomogeneousToCartesian()
						if pp.Z > -1 && pp.Z < 1 {
							acc.Blot(pp.X, pp.Y, c)
						}
					}
				}
			}
		}

		a := rndSphere(1)
		for i := 0; i < 10000; i++ {
			b := rndSphere(0.07)
			b = a.Add(b)
			b = b.Normalize()

			//mmm := mv.Mult(mm.M44())
			//mmm := mm.M44().Inverse().Mult(mv)
			//mmm := mm.M44().Mult(mv.Inverse()).Inverse()

			acc.BlotLine(mm.M44(), mv, cam, vector.Line{Start: a, End: b})

			a = b
		}

		/*	for x := -1.0; x < 1.0; x += 0.1 {
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{x, -1, -1},
				End: vector.V3{x, 1, -1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{x, -1, -1},
				End: vector.V3{x, -1, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, x, -1},
				End: vector.V3{1, x, -1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, x, -1},
				End: vector.V3{-1, x, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, -1, x},
				End: vector.V3{1, -1, x},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, -1, x},
				End: vector.V3{-1, 1, x},
			})

			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{x, -1, 1},
				End: vector.V3{x, 1, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{x, 1, -1},
				End: vector.V3{x, 1, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, x, 1},
				End: vector.V3{1, x, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{1, x, -1},
				End: vector.V3{1, x, 1},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{-1, 1, x},
				End: vector.V3{1, 1, x},
			})
			acc.BlotLine(mv, cam, vector.Line{
				Start: vector.V3{1, -1, x},
				End: vector.V3{1, 1, x},
			})
		}*/

		acc.Bg()

		img := acc.ToNRGBA()
		w, err := os.Create(fmt.Sprintf("frame%08d.png", frame))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.Close()
		if err := png.Encode(w, img); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		frame++
	}

}