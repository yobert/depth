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

const (
	m = 0.7
	//m = 0.1
	//m = 0.0
	e = 3.0
	//	m = 0.0
	//	e = 0.0
	f = 1.3

	//	m = 0

	motion = 100
	speed  = 0.01
	save   = false
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

func (img Fimg) BlotPoint(lr *rand.Rand, cam vector.Camera, p vector.V3, c Fcolor) {
	dist := p.Dist(cam.Position)

	// focal tricks
	r := m * math.Pow(math.Abs(f-dist), e)
	p = p.Add(vector.RandV3(lr).Scale(r))
	dist = p.Dist(cam.Position)

	pp := cam.ModelViewProjection.MultV4(p.CartesianToHomogeneous()).HomogeneousToCartesian()

	//	p4 := vector.V4{p.X, p.Y, p.Z, 1}
	//	p4 = cam.ModelViewProjection.MultV4(p4)
	//	pp := p4.HomogeneousToCartesian()

	// out of the clipping plane?
	if pp.Z < -1 || pp.Z > 1 {
		return
	}

	x := pp.X
	y := pp.Y

	w, h := img.Rect.Dx(), img.Rect.Dy()

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

	ix += img.Rect.Min.X
	iy += img.Rect.Min.Y

	a := c.A

	//radius := 3.0
	strength := 1.0

	//atten := clamp(1.0 - (dist / radius))
	//atten *= atten

	atten := 1.0 / (dist * dist)

	atten *= strength

	//fmt.Println(dist)
	//atten := pp.Z //* -0.5 + 1.0
	//atten := math.Pow(2, -dist)
	//fmt.Println(atten)

	//a *= 1.0 / math.Pow(atten, -2)
	//fmt.Println(atten)

	//atten = 1

	a *= atten

	_ = rxb
	_ = ryb
	c.A = a
	//c.A = 1
	img.Add(ix, iy, c)

	// lameo antialiasing
	//c.A = a * rxb * ryb
	//img.Add(ix+0, iy+0, c)
	//c.A = a * rxa * ryb
	//img.Add(ix+1, iy+0, c)
	//c.A = a * rxb * rya
	//img.Add(ix+0, iy+1, c)
	//c.A = a * rxa * rya
	//img.Add(ix+1, iy+1, c)
}

func (f Fimg) Add(x, y int, c Fcolor) {
	cc := f.Get(x, y)
	cc.R += c.R * c.A
	cc.G += c.G * c.A
	cc.B += c.B * c.A
	f.Set(x, y, cc)
}

func (img Fimg) BlotLine(lr *rand.Rand, cam vector.Camera, line vector.Line, samples int, c Fcolor) {
	for i := 0; i < samples; i++ {
		v := line.Lerp(lr.Float64())
		//v := line.Lerp(float64(i) / float64(samples))

		d := cam.Position.Dist(v)

		//a := 1.0 / (math.Pow(d, 2))
		//a *= 0.1
		//a := 0.005
		//		a := 1.0 / float64(samples)

		r := m * math.Pow(math.Abs(f-d), e)
		w := v.Add(vector.RandV3(lr).Scale(r))
		//		w := v

		//a := 0.001
		//		a := 0.1
		a := 1.0

		cc := c
		cc.A *= a
		img.BlotPoint(lr, cam, w, cc)
	}
}

func slide(v vector.V3) vector.V3 {
	return vector.V3{v.Y, v.Z, v.X}
}
func flip(v vector.V3) vector.V3 {
	return vector.V3{v.X, v.Z, v.Y}
}
func render(input inputType) *Fimg {
	//mm := vector.RotateAxisM33(vector.V3{0, -1, 0}, input.Rot)
	//mm := vector.RotateAxisM33(vector.V3{-2, -1, 0}.Normalize(), input.Rot).M44()

	mm := vector.RotateAxisM33(vector.V3{-2, -1, 0}.Normalize(), input.Rot).M44()
	//mm[14] += math.Cos(float64(input.Rot * 10))

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
		//1.5, -2)
		1.5, -3)

	cam.SetupModelView()

	cam.ModelViewProjection = cam.ModelView.MultX(cam.Projection)

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

	white := Fcolor{1, 1, 1, 1}
	samplefactor := 100

	/*	alph := 1.0

		for i := 0; i < 100000; i++ {
			a := vector.RandV3(lr).Normalize()
			a = mm.MultV3(a)
			acc.BlotPoint(lr, cam, a, Fcolor{1, 1, 1, alph})
		}*/

	cs := 2 * samplefactor
	_ = mm
	_ = white
	_ = cs

	for v := -1.0; v < 1.0; v += 0.1 {

		a := vector.V3{v, -1, -1}
		b := vector.V3{v, 1, -1}

		for iii := 0; iii < 2; iii++ {
			for ii := 0; ii < 2; ii++ {
				for i := 0; i < 3; i++ {
					acc.BlotLine(lr, cam, vector.Line{mm.MultV3(a), mm.MultV3(b)}, cs, white)
					a = slide(a)
					b = slide(b)
				}
				a = flip(a)
				b = flip(b)
			}
			a.Z *= -1
			b.Z *= -1
		}
	}

	/*	acc.BlotLine(lr, cam, vector.Line{corners[0], corners[1]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[2], corners[3]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[4], corners[5]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[6], corners[7]}, cs, white)

		acc.BlotLine(lr, cam, vector.Line{corners[0], corners[2]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[1], corners[3]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[4], corners[6]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[5], corners[7]}, cs, white)

		acc.BlotLine(lr, cam, vector.Line{corners[0], corners[4]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[1], corners[5]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[2], corners[6]}, cs, white)
		acc.BlotLine(lr, cam, vector.Line{corners[3], corners[7]}, cs, white)*/

	//acc.Shine()

	//	acc.BlotPoint(lr, cam, vector.V3{}, white)

	acc.Bg()
	acc.Clamp()

	return acc

	//	acc.Gamma(2.2)

	//	img := acc.ToNRGBA()

	//	if input.Frame < 0 {
	/*		w, err := os.Create(fmt.Sprintf("frame%08d.png", input.Frame))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer w.Close()
			if err := png.Encode(w, img); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}*/
	//	}

	//	return img
}

func main() {
	width := 1920  // 2
	height := 1080 // 2

	za := app.New()
	w := za.NewWindow("main")

	showchan := make(chan *image.NRGBA)
	inputs := make(chan inputType)

	go func() {
		frame := 0
		//for {
		for rot := vector.Radian(0); rot < 2*math.Pi*10; rot += (speed / motion) {
			inputs <- inputType{Width: width, Height: height, Rot: rot, Frame: frame}
			frame++
		}
		//}
		close(inputs)
	}()

	go func() {
		results := make(chan chan *Fimg, 8)

		go func() {
			for input := range inputs {
				rc := make(chan *Fimg)
				results <- rc
				go func(input inputType, rc chan *Fimg) {
					rc <- render(input)
					close(rc)
				}(input, rc)
			}
			close(results)
		}()

		var avg *Fimg

		count := 0

		for rc := range results {
			r := <-rc

			if avg == nil {
				avg = r
			} else {
				for i, p := range r.Pix {
					avg.Pix[i] += p
				}
			}
			count++

			if count == motion {
				if motion > 1 {
					for i, p := range avg.Pix {
						avg.Pix[i] = p / motion
					}
				}

				avg.Gamma(2.2)
				img := avg.ToNRGBA()
				showchan <- img
				avg = nil
				count = 0
			}
		}

		close(showchan)
	}()

	baseimg := image.NRGBA{}

	ica := canvas.NewRasterFromImage(&baseimg)
	ica.SetMinSize(fyne.NewSize(width, height))

	vb := widget.NewVBox(ica)

	go func() {

		frame := 0
		for img := range showchan {
			//time.Sleep(time.Second / 120)
			//time.Sleep(time.Second / 60)
			//			time.Sleep(time.Second / 24)
			_ = time.Sleep
			//time.Sleep(1)

			if save {
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
			}
			frame++

			baseimg = *img
			widget.Refresh(vb)
		}
	}()

	w.SetContent(vb)
	w.Show()
	za.Run()
	return
}
