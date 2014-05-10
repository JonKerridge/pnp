package pnp

import (
"testing"
)

func TestNumbers (t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	go Prefix(c, a, 0)
	go Copy2(a, b, d)
	go Plus1(b, c)
	var sum, check, i  int = 0, 0, 0
	for i < 10 {
		sum = sum + i
		check = check + <-d
		i = i + 1
	}
	if check != sum {
		t.Errorf ("Number funcs - check: %v sum: %v", check, sum )
	}
}

func TestIntegrate (t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	e := make(chan int)
	f := make(chan int)
	g := make(chan int)
	h := make(chan int)
	go Prefix(c, a, 0)
	go Copy2(a, b, d)
	go Plus1(b, c)
	go Plus (d, g, e)
	go Prefix(f, g, 0)
	go Copy2(e, f, h)
	var sum, check, i  int = 0, 0, 0
	for i < 10 {
		sum = sum + i
		check =  <-h
		i = i + 1
	}
	if check != sum {
		t.Errorf ("Integrate Funcs - check: %v sum: %v", check, sum )
	}
}

func TestIntegrateP (t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	e := make(chan int)
	f := make(chan int)
	g := make(chan int)
	h := make(chan int)
	go Prefix(c, a, 0)
	go Copy2(a, b, d)
	go Plus1(b, c)
	go PlusP (d, g, e)  // does not work with PlusP but does with Plus
	go Prefix(f, g, 0)
	go Copy2(e, f, h)
	var sum, check, i  int = 0, 0, 0
	for i < 10 {
		sum = sum + i
		check =  <-h
		i = i + 1
	}
	if check != sum {
		t.Errorf ("IntegrateP Funcs - check: %v sum: %v", check, sum )
	}
}

func TestReverseIntegrate (t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	e := make(chan int)
	f := make(chan int)
	g := make(chan int)
	h := make(chan int)
	i := make(chan int)
	j := make(chan int)
	k := make(chan int)
	l := make(chan int)
	go Prefix(c, a, 0)
	go Copy2(a, b, d)
	go Plus1(b, c)
	go Plus (d, g, e)
	go Prefix(f, g, 0)
	go Copy2(e, f, h)
	go Minus (i, k, l)
	go Prefix(j, k, 0)
	go Copy2(h, i, j)
	var check, n  int = 0, 0
	for n < 10 {
		check =  <-l
		n = n + 1
	}
	if check != n-1 {
		t.Errorf ("Reverse Integrate Funcs - check: %v i: %v", check, n )
	}
}

func TestSquares (t *testing.T) {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	e := make(chan int)
	f := make(chan int)
	g := make(chan int)
	h := make(chan int)
	i := make(chan int)
	j := make(chan int)
	k := make(chan int)
	l := make(chan int)
	go Prefix(c, a, 0)
	go Copy2(a, b, d)
	go Plus1(b, c)
	go Plus (d, g, e)
	go Prefix(f, g, 0)
	go Copy2(e, f, h)
	go Plus (i, k, l)
	go Tail(j, k)
	go Copy2(h, i, j)
	var check, n  int = 0, 0
	for n < 10 {
		check =  <-l
		n = n + 1
	}
	var square int = n * n
	if check != square {
		t.Errorf ("Square Funcs - check: %v square: %v", check, square )
	}
}

