package main

import (
	"math"
	"math/cmplx"
	"runtime"
	"fmt"
	"github.com/mjibson/go-dsp/dsputils"
	"github.com/mjibson/go-dsp/fft"
)

const (
	sqrt2_2 = math.Sqrt2 / 2
)

type fftTest struct {
	in  []float64
	out []complex128
}

var fftTests = []fftTest{
	// impulse responses
	{
		[]float64{1},
		[]complex128{complex(1, 0)},
	},
	{
		[]float64{1, 0},
		[]complex128{complex(1, 0), complex(1, 0)},
	},
	{
		[]float64{1, 0, 0, 0},
		[]complex128{complex(1, 0), complex(1, 0), complex(1, 0), complex(1, 0)},
	},
	{
		[]float64{1, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},

	// shifted impulse response
	{
		[]float64{0, 1},
		[]complex128{complex(1, 0), complex(-1, 0)},
	},
	{
		[]float64{0, 1, 0, 0},
		[]complex128{complex(1, 0), complex(0, -1), complex(-1, 0), complex(0, 1)},
	},
	{
		[]float64{0, 1, 0, 0, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(sqrt2_2, -sqrt2_2),
			complex(0, -1),
			complex(-sqrt2_2, -sqrt2_2),
			complex(-1, 0),
			complex(-sqrt2_2, sqrt2_2),
			complex(0, 1),
			complex(sqrt2_2, sqrt2_2)},
	},

	// other
	{
		[]float64{1, 2, 3, 4},
		[]complex128{
			complex(10, 0),
			complex(-2, 2),
			complex(-2, 0),
			complex(-2, -2)},
	},
	{
		[]float64{1, 3, 5, 7},
		[]complex128{
			complex(16, 0),
			complex(-4, 4),
			complex(-4, 0),
			complex(-4, -4)},
	},
	{
		[]float64{1, 2, 3, 4, 5, 6, 7, 8},
		[]complex128{
			complex(36, 0),
			complex(-4, 9.65685425),
			complex(-4, 4),
			complex(-4, 1.65685425),
			complex(-4, 0),
			complex(-4, -1.65685425),
			complex(-4, -4),
			complex(-4, -9.65685425)},
	},

	// non power of 2 lengths
	{
		[]float64{1, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},
	{
		[]float64{1, 2, 3},
		[]complex128{
			complex(6, 0),
			complex(-1.5, 0.8660254),
			complex(-1.5, -0.8660254)},
	},
	{
		[]float64{1, 1, 1},
		[]complex128{
			complex(3, 0),
			complex(0, 0),
			complex(0, 0)},
	},
}

func TestFFT() {
	for _, ft := range fftTests {
		v := fft.FFTReal(ft.in)
		if !dsputils.PrettyCloseC(v, ft.out) {
			//t.Error("FFT error\ninput:", ft.in, "\noutput:", v, "\nexpected:", ft.out)
		}

		vi := fft.IFFT(ft.out)
		if !dsputils.PrettyCloseC(vi, dsputils.ToComplex(ft.in)) {
			//t.Error("IFFT error\ninput:", ft.out, "\noutput:", vi, "\nexpected:", dsputils.ToComplex(ft.in))
		}
	}
}

func main() {
	TestFFT()
}
