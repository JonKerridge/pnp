package pnp

import (
	"testing"
)

func TestNumbersFunc(t *testing.T) {
	n2p := make(chan int)
	go Numbers(n2p, 1)
	var sum, check, i int = 0, 0, 1
	for i < 10 {
		sum = sum + i
		check = check + <-n2p
		i = i + 1
	}
	if check != sum {
		t.Errorf("Numbers - check: %v sum: %v", check, sum)
	}
}

func TestIntegrateFunc (t *testing.T) {
	n2i := make(chan int)
	i2p := make(chan int)
	go Numbers(n2i, 0)
	go Integrate(n2i, i2p)
	var sum, check, i int = 0, 0, 0
	for i < 10 {
		sum = sum + i
		check =  <- i2p
		i = i + 1
	}
	check =  <- i2p
	if check != sum + i{
		t.Errorf("Integrate - check: %v sum: %v", check, sum)
	}
}


func TestSquaresFunc (t *testing.T) {
	n2i := make(chan int)
	i2p := make(chan int)
	out := make(chan int)
	go Numbers (n2i, 0 )
	go Integrate (n2i, i2p)
	go Pairs (i2p, out)	
	var check, i int = 0, 0
	for i < 9 {
		check =  <- out
		i = i + 1
	}
	if check != 81 {
		t.Errorf("Squares - square: %v expected: %v", check, 81)
	}

}
func TestReversIntegrate (t *testing.T) {
	n2i := make(chan int)
	i2r := make(chan int)
	out := make(chan int)
	go Numbers (n2i, 0 )
	go Integrate (n2i, i2r)
	go ReverseIntegrate (i2r, out)	
	var check, i int = 0, 0
	for i < 9 {
		check =  <- out
		i = i + 1
	}
	check =  <- out
	if check != 9 {
		t.Errorf("Reverse Integrate - value: %v expected: %v", check, 9)
	}

}
