package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"os"

	"github.com/zephyrtronium/xirho/encoding"
)

func main() {
	var (
		decode bool
	)
	flag.BoolVar(&decode, "d", false, "decode rather than encode")
	flag.Parse()

	if decode {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		doDecode(string(b))
	} else {
		doEncode(os.Stdin)
	}
}

func doDecode(p string) {
	palette, err := encoding.DecodePalette(p)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range palette {
		r, g, b, a := scale(c.RGBA())
		fmt.Printf("%f\t%f\t%f\t%f\n", r, g, b, a)
	}
}

func scale(r, g, b, a uint32) (rf, gf, bf, af float32) {
	rf = float32(r) / 0xffff
	gf = float32(g) / 0xffff
	bf = float32(b) / 0xffff
	af = float32(a) / 0xffff
	return
}

func doEncode(in io.Reader) {
	var palette color.Palette
	for {
		var r, g, b, a float32
		_, err := fmt.Fscanf(in, "%g %g %g %g\n", &r, &g, &b, &a)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		palette = append(palette, unscale(r, g, b, a))
	}
	fmt.Println(encoding.EncodePalette(palette))
}

func unscale(r, g, b, a float32) color.RGBA64 {
	return color.RGBA64{
		R: uint16(r * 0xffff),
		G: uint16(g * 0xffff),
		B: uint16(b * 0xffff),
		A: uint16(a * 0xffff),
	}
}
