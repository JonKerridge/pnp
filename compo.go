package pnp

/*
Numbers generates a sequence of increasing integers 
that it sends on its out channel starting from initialValue.
*/
func Numbers ( out chan int,
	initialValue int ){
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go Prefix(c, a, initialValue)
	go Copy2(a, b, out)
	go Plus1(b, c)
}	

/*
Integrate receives values on its in channel and sends
the running sum of the values on its out channel
*/
func Integrate ( in chan int,
	out chan int){
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go Prefix(b, c, 0)
	go Plus(in, c, a)
	go Copy2 ( a, b, out )
}

/*
Pairs receives a stream of values on its in channel and sends the
sum of consecutive pairs on its out channel.
Thus the input stream 0,1,2,3,4 
would result in the output
1,3,5,7
*/
func Pairs ( in chan int,
	out chan int){
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go Copy2 (in, a, b)
	go Tail (b, c)
	go Plus (a, c, out)	
}	

/*
Squares is a network of Numbers Integrate and Pairs that
sends the sqaures of the integers to its out channel
*/
func Squares ( out chan int ){
	n2i := make(chan int)
	i2p := make(chan int)
	go Numbers (n2i, 0 )
	go Integrate (n2i, i2p)
	go Pairs (i2p, out)	
}	

/*
ReverseIntegrate undertakes the opposite action of Integrate.
*/
func ReverseIntegrate ( in chan int,
	out chan int){
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	go Copy2 (in, a, b)
	go Prefix (b, c, 0)
	go Minus (a, c, out)
}