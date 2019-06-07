package main

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand"
	"os"
	"time"

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

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func (f Fimg) BlotPoint(lr *rand.Rand, pos vector.V3, mv vector.M44, p vector.V3, c Fcolor) {
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

	ix := int(x)
	iy := int(y)

	rxa := x - float64(ix)
	rya := y - float64(iy)

	rxb := 1.0 - rxa
	ryb := 1.0 - rya

	ix += f.Rect.Min.X
	iy += f.Rect.Min.Y

	a := c.A

	dist := p.Dist(pos)

	radius := 3.0
	strength := 3.0

	atten := clamp(1.0 - (dist / radius))
	atten *= atten

	atten *= strength

	//fmt.Println(dist)
	//atten := pp.Z //* -0.5 + 1.0
	//atten := math.Pow(2, -dist)
	//fmt.Println(atten)

	//a *= 1.0 / math.Pow(atten, -2)
	//fmt.Println(atten)
	a *= atten

	c.A = a * rxb * ryb
	f.Add(ix+0, iy+0, c)

	c.A = a * rxa * ryb
	f.Add(ix+1, iy+0, c)

	c.A = a * rxb * rya
	f.Add(ix+0, iy+1, c)

	c.A = a * rxa * rya
	f.Add(ix+1, iy+1, c)
}

func (f Fimg) Add(x, y int, c Fcolor) {
	cc := f.Get(x, y)
	cc.R += c.R * c.A
	cc.G += c.G * c.A
	cc.B += c.B * c.A
	f.Set(x, y, cc)
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
		img.BlotPoint(lr, cam.Position, mv, w, c)
	}
}

func render(input inputType) interface{} {
	//mm := vector.RotateAxisM33(vector.V3{0, -1, 0}, input.Rot)
	mm := vector.RotateAxisM33(vector.V3{-1, -1, 0}, input.Rot)

	lr := rand.New(rand.NewSource(1).(rand.Source64))
	lr.Seed(666)

	sx := input.Width / 2
	sy := input.Height / 2

	acc := NewFimg(image.Rect(-sx, -sy, sx, sy))

	//acc.Set(0, 0, Fcolor{1, 1, 1, 1})

	cam := vector.Camera{
		Width:  float64(input.Width),
		Height: float64(input.Height),

		YFov: 100,

		Near: 0.1,
		Far:  20,

		Position: vector.V3{0, 0, 1.2},
		RotAxis:  vector.Euler{
			//			X: 0.01,
		},
	}

	//cam.SetupViewProjection()
	//cam.View = vector.Frustum{
	//	-float64(sy), float64(sx),
	//	float64(sy), -float64(sx),
	//	-1, 1,
	//}
	//cam.Projection = cam.View.M44()

	//cam.Projection = vector.Ortho(
	//	-float64(sx), float64(sx),
	//	-float64(sy), float64(sy),
	//	-1, 1)
	x_ratio := float64(sx) / float64(sy)
	cam.Projection = vector.Ortho(
		-x_ratio, x_ratio,
		-1, 1,
		1.5, -2)

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

	_ = mm

	/*	a := vector.RandV3(lr)
		for i := 0; i < 1000; i++ {
			b := vector.RandV3(lr).Scale(0.7)
			b = a.Add(b)
			b = b.Normalize()

			//mmm := mv.Mult(mm.M44())
			//mmm := mm.M44().Inverse().Mult(mv)
			//mmm := mm.M44().Mult(mv.Inverse()).Inverse()

			acc.BlotLine(lr, mm.M44(), mv, cam, vector.Line{Start: a, End: b})

			a = b
		}*/

	for i := 0; i < 100; i++ {
		b := vector.RandV3(lr)
		b = b.Normalize()

		b = mm.MultV3(b)

		acc.BlotPoint(lr, cam.Position, mv, b, Fcolor{1, 1, 1, 1})
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

	//acc.Shine()

	acc.Bg()

	img := acc.ToNRGBA()

	if input.Frame < 0 {
		w, err := os.Create(fmt.Sprintf("frame%08d.png", input.Frame))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.Close()
		if err := png.Encode(w, img); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return img
}

func main() {
	width := 1920
	height := 1080

	za := app.New()
	w := za.NewWindow("main")

	showchan := make(chan interface{})
	inputs := make(chan inputType)

	go func() {
		frame := 0
		for {
			for rot := vector.Radian(0); rot < 2*math.Pi; rot += 0.01 {
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
