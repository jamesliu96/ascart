package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var table = func() (t []byte) {
	T := [][]byte{
		{' '},
		{'`'},
		{'\''},
		{'-', '.', '^'},
		{'"', ',', '_'},
		{'~'},
		{'!'},
		{'(', ')', '+', '/', '<', '>', '\\', '|'},
		{':', 'r'},
		{';', 'v'},
		{'=', 'i', 'x'},
		{'?', 'T', 'Y'},
		{'7', '*', 'J', 'L', 'c', 'j', 'l', 't'},
		{'1', 'f', 'n', 'u', 'w'},
		{'k', 'o', 's', 'z', '{', '}'},
		{'$', '%', '&', 'C', 'F', 'I', '[', ']', 'm'},
		{'0', 'K', 'V', 'X', 'h'},
		{'4', 'a', 'e', 'p', 'q'},
		{'2', '3', '#', 'P', 'S', 'U', 'Z', 'y'},
		{'6', '9', 'b', 'd'},
		{'A', 'E', 'G', 'H', 'N', 'O', 'g'},
		{'5'},
		{'8', 'D', 'R'},
		{'Q'},
		{'B', 'M', 'W'},
		{'@'},
	}
	t = make([]byte, len(T))
	for i, c := range T {
		t[i] = c[rand.Intn(len(c))]
	}
	return
}()

var lentable = len(table)

func getByte(n uint8) byte {
	return table[int(float64(n)/math.MaxUint8*float64(lentable-1))]
}

func drawByte(dst draw.Image, c color.Color, x, y int, s byte) {
	(&font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, basicfont.Face7x13.Ascent+y),
	}).DrawBytes([]byte{s})
}

func main() {
	src, _, err := image.Decode(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	b := src.Bounds()
	dx, dy := b.Dx(), b.Dy()
	dst := image.NewRGBA(b)
	w, h := basicfont.Face7x13.Advance, basicfont.Face7x13.Height
	for oy := 0; oy < dy; oy += h {
		for ox := 0; ox < dx; ox += w {
			sum, cnt := 0, 0
			for y := oy; y < oy+h && y < dy; y++ {
				for x := ox; x < ox+w && x < dx; x++ {
					sum += int(color.GrayModel.Convert(src.At(x, y)).(color.Gray).Y)
					cnt++
				}
			}
			b := getByte(uint8(sum / cnt))
			drawByte(dst, color.White, ox, oy, b)
			fmt.Fprintf(os.Stderr, "%c", b)
		}
		fmt.Fprintln(os.Stderr)
	}
	if err = png.Encode(os.Stdout, dst); err != nil {
		log.Fatalln(err)
	}
}
