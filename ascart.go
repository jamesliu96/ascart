package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

var t = [][]byte{
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

var stable = func() (tt []byte) {
	tt = make([]byte, len(t))
	for i, c := range t {
		tt[i] = c[0]
	}
	return
}()

var rtable = func() (tt []byte) {
	tt = make([]byte, len(t))
	for i, c := range t {
		tt[i] = c[rand.Intn(len(c))]
	}
	return
}()

var lentable = len(t)

func getByte(n uint16, r bool) byte {
	x := int(float64(n) / math.MaxUint16 * float64(lentable-1))
	if r {
		return rtable[x]
	}
	return stable[x]
}

func drawByte(dst draw.Image, c color.Color, x, y int, s byte) {
	if s == ' ' {
		return
	}
	(&font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(c),
		Face: basicfont.Face7x13,
		Dot:  fixed.P(x, basicfont.Face7x13.Ascent+y),
	}).DrawBytes([]byte{s})
}

func ascart(r io.Reader, w io.Writer, colored bool, random bool) (err error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return
	}
	b := src.Bounds()
	dx, dy := b.Dx(), b.Dy()
	dst := image.NewRGBA64(b)
	W, H := basicfont.Face7x13.Advance, basicfont.Face7x13.Height
	for oy := 0; oy < dy; oy += H {
		for ox := 0; ox < dx; ox += W {
			sumY, sumR, sumG, sumB, cnt := 0, 0, 0, 0, 0
			for y := oy; y < oy+H && y < dy; y++ {
				for x := ox; x < ox+W && x < dx; x++ {
					c := src.At(x, y)
					y := color.Gray16Model.Convert(c).(color.Gray16)
					sumY += int(y.Y)
					if colored {
						nrgba := color.NRGBA64Model.Convert(c).(color.NRGBA64)
						sumR += int(nrgba.R)
						sumG += int(nrgba.G)
						sumB += int(nrgba.B)
					}
					cnt++
				}
			}
			var dstC color.Color = color.White
			if colored {
				dstC = color.RGBA64{uint16(sumR / cnt), uint16(sumG / cnt), uint16(sumB / cnt), math.MaxUint16}
			}
			drawByte(dst, dstC, ox, oy, getByte(uint16(sumY/cnt), random))
		}
	}
	err = png.Encode(w, dst)
	return
}

var (
	fColored = flag.Bool("c", false, "colored output")
	fRandom  = flag.Bool("r", false, "random output")
)

func main() {
	flag.Parse()
	if err := ascart(os.Stdin, os.Stdout, *fColored, *fRandom); err != nil {
		log.Fatalln(err)
	}
}
