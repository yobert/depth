package main

import (
	"fmt"
	"image"
	"time"
	//	"image/png"
	"math"
	"math/rand"
	//	"os"

	"depth/vector"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

type inputType struct {
	Width, Height int
	Rot           vector.Radian
	Frame         int
}

func (f Fimg) BlotPoint(lr *rand.Rand, mv vector.M44, p vector.V3, c Fcolor) {
	p4 := vector.V4{p.X, p.Y, p.Z, 1}
	p4 = mv.MultV4(p4)
	pp := p4.HomogeneousToCartesian()

	// out of the clipping plane?
	if pp.Z < -1 || pp.Z > 1 {
		return
	}

	x := pp.X
	y := pp.Y

	w, h := f.Rect.Dx(), f.Rect.Dy()

	x += 1
	y += 1

	x *= 0.5
	y *= 0.5

	x *= float64(w)
	y *= float64(h)

	antialias := 1
	for i := 0; i < antialias*4; i++ {
		ix := f.Rect.Min.X + int(x+(rand.Float64()*2.0-1.0))
		iy := f.Rect.Min.Y + int(y+(rand.Float64()*2.0-1.0))

		oldc := f.Get(ix, iy)

		newc := Fcolor{
			c.R*c.A/float64(antialias) + oldc.R,
			c.G*c.A/float64(antialias) + oldc.G,
			c.B*c.A/float64(antialias) + oldc.B,
			0,
		}

		f.Set(ix, iy, newc)
	}
}

func (img Fimg) BlotLine(lr *rand.Rand, mm, mv vector.M44, cam vector.Camera, ray vector.Line) {
	//m := 0.3
	//e := 3.0
	m := 0.0
	e := 0.0
	f := 1.3

	for i := 0; i < int(ray.Start.Dist(ray.End)*1000.0); i++ {
		v := ray.Lerp(lr.Float64())

		v = mm.MultV3(v)

		d := cam.Position.Dist(v)

		//a := 1.0 / (math.Pow(d, 2))
		//a *= 0.1
		a := 0.005
		//a := 1.0

		r := m * math.Pow(math.Abs(f-d), e)
		w := v.Add(vector.RandV3(lr).Scale(r))

		c := Fcolor{1, 1, 1, a}
		img.BlotPoint(lr, mv, w, c)
	}
}

func render(input inputType) interface{} {
	mm := vector.RotateAxisM33(vector.V3{0, -1, 0}, input.Rot)

	lr := rand.New(rand.NewSource(1).(rand.Source64))
	lr.Seed(666)

	sx := input.Width / 2
	sy := input.Height / 2

	acc := NewFimg(image.Rect(-sx, -sy, sx-1, sy-1))

	cam := vector.Camera{
		Width:  float64(input.Width),
		Height: float64(input.Height),

		YFov: 120,
		Near: 0.1,
		Far:  100,

		Position: vector.V3{0, 0, 1.2},
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

	a := vector.RandV3(lr)
	for i := 0; i < 1000; i++ {
		b := vector.RandV3(lr).Scale(0.7)
		b = a.Add(b)
		b = b.Normalize()

		//mmm := mv.Mult(mm.M44())
		//mmm := mm.M44().Inverse().Mult(mv)
		//mmm := mm.M44().Mult(mv.Inverse()).Inverse()

		acc.BlotLine(lr, mm.M44(), mv, cam, vector.Line{Start: a, End: b})

		a = b
	}

	/*	for i := 0; i < 1000; i++ {
		b := vector.RandV3(lr)
		b = b.Normalize()

		b = mm.MultV3(b)

		acc.BlotPoint(lr, mv, b, Fcolor{1, 1, 1, 1})
	}*/

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

	/*		w, err := os.Create(fmt.Sprintf("frame%08d.png", frame))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer w.Close()
			if err := png.Encode(w, img); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		}*/

	return img
}

func main() {
	width := 640
	height := 480

	za := app.New()
	w := za.NewWindow("main")

	showchan := make(chan interface{})
	inputs := make(chan inputType)

	go func() {
		frame := 0
		for {
			for rot := vector.Radian(0); rot < 2*math.Pi; rot += 0.001 {
				inputs <- inputType{Width: width, Height: height, Rot: rot, Frame: frame}
				frame++
			}
		}
		close(inputs)
	}()

	go func() {
		results := make(chan chan interface{}, 4)

		go func() {
			for input := range inputs {
				rc := make(chan interface{})
				results <- rc
				go func(input inputType, rc chan interface{}) {
					rc <- render(input)
					close(rc)
				}(input, rc)
			}
			close(results)
		}()

		for rc := range results {
			r := <-rc
			showchan <- r
		}
		close(showchan)
	}()

	baseimg := image.NRGBA{}

	ica := canvas.NewRasterFromImage(&baseimg)
	ica.SetMinSize(fyne.NewSize(width, height))

	vb := widget.NewVBox(ica)

	go func() {

		for res := range showchan {
			//time.Sleep(time.Second / 120)
			time.Sleep(time.Second / 60)
			//time.Sleep(1)
			img, ok := res.(*image.NRGBA)
			if !ok {
				fmt.Printf("not ok: %T\n", res)
				continue
			}

			baseimg = *img
			widget.Refresh(vb)
		}
	}()

	w.SetContent(vb)
	w.Show()
	za.Run()
	return
}
