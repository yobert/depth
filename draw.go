package main

import (
	"image"
	"image/color"
)

func hline(img *image.NRGBA, x1, y, x2 int, col color.Color) {
	for ; x1 <= x2; x1++ {
		img.Set(x1, y, col)
	}
}

func vline(img *image.NRGBA, x, y1, y2 int, col color.Color) {
	for ; y1 <= y2; y1++ {
		img.Set(x, y1, col)
	}
}

func box(img *image.NRGBA, x1, y1, x2, y2 int, col color.Color) {
	for y := y1; y <= y2; y++ {
		hline(img, x1, y, x2, col)
	}
}

/*func rect(x1, y1, x2, y2 int) {
    HLine(x1, y1, x2)
    HLine(x1, y2, x2)
    VLine(x1, y1, y2)
    VLine(x2, y1, y2)
}*/
