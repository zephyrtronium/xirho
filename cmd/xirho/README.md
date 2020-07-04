# xirho

The xirho command implements a basic renderer using xirho.

## Installation

To install the xirho command, simply run `go get github.com/zephyrtronium/xirho/cmd/xirho` with a Go version of at least 1.14.

## Usage

By default, xirho reads a system from stdin and writes the resulting 1024x1024 PNG image to stdout. There are a number of command-line options to control the render, e.g. `-width` and `-height` to set the image size, `-dur` to set the amount of time to spend rendering, and so on. Most of the images in xirho/img were rendered using an invocation like:

`xirho -gamma 2.2 -osa 6 -png "test.png" -dur 2m -width 1024 -height 576 <img/discjulian.json`

See `xirho -help` for more details.

Note that to use xirho, you need fractal parameters. See img/xirho for some simple examples, or try using an Apophysis flame file with the `-flame` option.
