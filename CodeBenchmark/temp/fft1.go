package main

import (
	_"math"
	//"os"
	//"strconv"
	"github.com/mjibson/go-dsp/fft"
	"fmt"
)



type fftTest struct {
	in  []float64
	out []complex128
}

func generateFftTest(n int) *[]float64{
	in := make([]float64,n)
	//out := make([]complex128,n)
	in = append(in,1)
	//out = append(out,complex(1,0))
	for i:=0 ; i< n-1 ; i++{
		in = append(in,0)
		//out = append(out,complex(1,0))
	}
	return &in
	//fftTest := &fftTest{in:in,out:out}
	//return fftTest
}


func main() {
	//_n := os.Getenv("FFT_SIZE")
	//n,_ := strconv.Atoi(_n)
	n := 32768
	ft := generateFftTest(n)
	fmt.Println(n)
	fft.FFTReal(*ft)
}

/*var fftTests = []fftTest{
	// impulse responses
	{
		[]float64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]complex128{
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0),
			complex(1, 0)},
	},
}*/
