package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events throughout code.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(arg ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(arg...)))
	t.out.Write([]byte("\n"))
}

// New - create a new Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// SayHello - say something - testing
func SayHello() {
	fmt.Println("Hi there!")
}
