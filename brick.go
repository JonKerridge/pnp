/*
Package pnp, (plug and play),  contains a number of funcs that are executed as
go routines.  They are generally very simple which makes the compositon of them into
more complex networks very easy.  The composition of them produces output
which at times may seem surprising.

Each func has a number of channels with a consistent naming convetion.
Channels that are used to receive values are named in.
Channels that are used to send values are named out.

Multiple channels are differentiated by appending a digit to the name.

Each func operates in a for {} loop and thus any network created has to
be terminated by the user unless the output of  a network is limited bt
a final phase that limits the output.

The functions Numbers and Squares show how many of these simple processes
can be used to create complex behaviour simply by composing them into networks.

The hello.go example shows how these can be made to work in a more complex example.

The exerciseQueue.go exxample shows the operation of Producer, Queue and Prompt.
*/
package pnp

import (
	"fmt"
	"sync"
	"time"
)

/*
Prefix outputs the initial value on the out channel and then sends
all values received on the in channel to the out channel unaltered.
*/
func Prefix(in chan int,
	out chan int,
	initial int) {
	out <- initial
	for {
		out <- (<-in)
	}
}

/*
ResetPrefix outputs the initial value on the out channel and then sends
all values received on the in channel to the out channel unaltered.
Any input on the channel reset overwrites the next input from the in
channel.
*/
func ResetPrefix(in chan int,
	out chan int,
	reset chan int,
	initial int) {
	var v, r int = 0, 0
	out <- initial
	for {
		select {
		case r = <-reset:
			<-in
			out <- r
		case v = <-in:
			out <- v
		}
	}
}

/*
Plus1 receives a value on its in channel and that value plus 1 is sent
to the out channel
*/
func Plus1(in chan int,
	out chan int) {
	for {
		out <- (<-in) + 1
	}
}

/*
RecvValue receives a single value from the in channel which it stores in
the variable valuePtr.  RecvValue does not contain a for {} loop;
it terminates once the value has been read from the in channel
See its use in PlusP and other functions
*/
func RecvValue(valuePtr *int,
	in chan int,
	wg *sync.WaitGroup) {

	*valuePtr = <-in
	wg.Done()
}

/*
SendValue sends the single value to the out channel.  SendValue does not
contain a for {} loop; it terminates once the value has been sent.
See its use in Copy2
*/
func SendValue(value int,
	out chan int) {
	out <- value
}

/*
Copy2 receives a value from its in channel, which it then sends to
each of the channels out1 and out2.  The sending of the values is achieved
using calls to go SendValue and thus happen concurrently
*/
func Copy2(in chan int,
	out1 chan int,
	out2 chan int) {
	for {
		v := <-in
		go SendValue(v, out1)
		go SendValue(v, out2)
	}
}

/*
CopyN receives a value on its in channel and then writes the same value
concurrently to all the channels in the slice of channels out.
*/
func CopyN(in chan int,
	out []chan int) {
	l := len(out)
	var i, v int = 0, 0
	for {
		i = 0
		v = <-in
		for i < l {
			go SendValue(v, out[i])
			i = i + 1
		}
	}
}

/*
Tail receives a value from the in channel which it then ignores, therafter
it sends all received values on the in channel to the out channel unaltered.
*/
func Tail(in chan int,
	out chan int) {
	v := <-in // receive and throw away
	for {
		v = <-in
		out <- v
	}
}

/*
PlusP receives two values one from each of its in channels in1 and in2
and it then sends the sum of the values to the out channel.
It uses go routine RecvValue to read the values
regardless of the order in which they arrive.

NOTE the sync/WaitGroup is required
because go routines have no concept of
coordinated termination as say happens in occam or JCSP
and also in Unix as in fork / join
*/
func PlusP(in1 chan int,
	in2 chan int,
	out chan int) {
	var v1, v2 int = 0, 0
	var wg sync.WaitGroup
	for {
		wg.Add(2)
		go RecvValue(&v1, in1, &wg)
		go RecvValue(&v2, in2, &wg)
		wg.Wait()
		out <- v1 + v2
	}
}

/*
Minus receives two values one from each of its in channels in1 and in2
concurrently and it then sends the difference of the values to the out channel.
The value from in2 is subtracted from that obtained from in1.
The contained select statement ensures that a pair of values are received,
one from each channel, regardless of the order in which they arrive.
*/
func Minus(in1 chan int,
	in2 chan int,
	out chan int) {
	var v1, v2 int = 0, 0
	var wg sync.WaitGroup
	for {
		wg.Add(2)
		go RecvValue(&v1, in1, &wg)
		go RecvValue(&v2, in2, &wg)
		wg.Wait()
		out <- v1 - v2
	}
}

/*
Plus receives two values one from each of its in channels in1 and in2
and it then sends the sum of the values to the out channel.
The contained select statement ensures that a pair of values are received,
one from each channel, regardless of the order in which they arrive.
*/
func Plus(in1 chan int,
	in2 chan int,
	out chan int) {
	var v1, v2 int = 0, 0
	for {
		//read in1 and in2 in pseudo parallel
		select {
		case v1 = <-in1:
			v2 = <-in2
		case v2 = <-in2:
			v1 = <-in1
		}
		out <- v1 + v2
	}
}

/*
Delay introduces a time delay into a network.

The current formulation does not work with go1.1.2.  An issue has been
raised to which the response was that the formulation should work with the next
release of the scheduler (15th August 2013)
*/
func Delay(in chan int,
	out chan int,
	seconds int64) {
	delay := time.Second * time.Duration(seconds)
	var v int = 0
	for {
		v = <-in
		time.Sleep(delay) //should work in go1.2
		out <- v
	}
}

/*
Consumer receives values from its in channel, which
are then printed on the console.

It currently does not work correctly due to scheduler problems see Delay.
*/
func Consumer(in chan int) {
	var v int = 0
	v = <-in
	for v >= 0 {
		fmt.Printf("%v\n", v)
		v = <-in
	}
}

/*
Producer sends the values 0 .. iterations upwards in steps of 1 to
the out channel. It then send the value -1 to the out channel which
can be used as a terminating value
*/
func Producer(out chan int,
	iterations int) {
	var v int = 0
	for v < iterations {
		out <- v
		v = v + 1
	}
	out <- -1
}

/*
Tabulate receives inputs from an array of in channels, which it tabulates
into a single line of text which is sent as a string to the out channel.

A value has to be received from every input channel before any output is
generated. This input is undertaken concurrently.
*/
func Tabulate(in []chan int,
	out chan string) {
	var s string
	var i int
	var wg sync.WaitGroup
	l := len(in)
	var values = make([]int, l)
	for {
		s = ""
		i = 0
		wg.Add(l)
		for i < l {
			go RecvValue(&values[i], in[i], &wg)
			i = i + 1
		}
		wg.Wait()
		i = 0
		for i < l {
			s = s + fmt.Sprintf("\t%v", values[i])
			i = i + 1
		}
		s = s + fmt.Sprintf("\n")
		out <- s
	}
}

/* ConvertIntStr is a process that converts and integer into its string
representation.  Int values are input from the in channel and the equivalent
string is output on the out channel
*/
func ConvertIntStr(in chan int,
	out chan string) {
	var v int = 0
	var s string
	for {
		v = <-in
		s = fmt.Sprintf(" %v, ", v)
		out <- s
	}
}

/*
Display takes a string input from the in channel and outputs its value
to stdout. It does this forever.
*/
func Display(in chan string) {
	var s string
	for {
		s = <-in
		fmt.Printf("%v\n", s)
	}
}

/*
Console is similar to Display in that string input from the in channel
are output to stdout.  Initially, a title string is output.  Only limit
lines are output
*/
func Console(in chan string,
	title string,
	limit int) {
	var s string
	var i int = 0
	fmt.Printf("%v\n", title)
	for i < limit {
		s = <-in
		fmt.Printf("%v\n", s)
		i = i + 1
	}
	fmt.Println("%v\n", "Console Output Finished")
}

/*
Queue implements a multi-element circular queue with elements slots.
It receives inputs to the Queue from its put channel.
On receiving a signal on iys get channel it outputs a value
from the Queue on the out channel.

The operational invariance of the Queue is maintained by ensuring that
the count of the number of filled slots follows the following:

count = 0 only receives on put permitted
count = elements only receives on get permitted
otherwise both receives on get and put are accepted.

A receive on get is immediately followed by a send of the next
value from the Queue to the out channel.


Queue operates as a component in a Client-Server network.

The Queue behaves as a pure server in that it responds to receives
on its input channels in finite time.

The example exerciseQueue shows how Queue can be used in
conjunction with Producer and Prompt, which both act as pure clients.
*/
func Queue(put chan int,
	get chan int,
	out chan int,
	elements int) {
	var front, rear, count int = 0, 0, 0
	var data = make([]int, elements)
	for {
		if count == 0 {
			// can only put elements into queue
			data[front] = <-put
			front = (front + 1) % elements
			count = count + 1
		} else if count == elements {
			// can only get data from queue
			<-get
			out <- data[rear]
			rear = (rear + 1) % elements
			count = count - 1
		} else {
			// can either do a get or put
			select {
			case data[front] = <-put:
				front = (front + 1) % elements
				count = count + 1
			case <-get:
				out <- data[rear]
				rear = (rear + 1) % elements
				count = count - 1
			}
		}
	}
}

/*
Prompt is a client that requests values from a server by sending
a signal on its get channel and then receiving a response on its recv channel.
The value is then output on the out channel.

It can be used in many networks where data is being streamed and can be used to
avoid deadlock and livelock or to provide a break in the data stream so that different parts
of the stream can run at different rates.
*/
func Prompt(get chan int,
	recv chan int,
	out chan int) {
	for {
		get <- 1
		out <- <-recv
	}
}
