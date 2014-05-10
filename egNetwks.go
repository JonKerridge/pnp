package pnp

/*
Example1 produces a well known sequence, which is left as an exercise for the
reader to work out how the result is achieved.
*/
func Example1 ( out chan int){
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	d := make(chan int)
	go Prefix(d, a, 0)
	go Prefix(c, d, 1)
	go Copy2(a, b, out)
	go Pairs(b, c)
}