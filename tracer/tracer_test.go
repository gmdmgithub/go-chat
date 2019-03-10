package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Fatal error. Buffer not creted - nil")
	} else {
		tracer.Trace("Tracer package created properly!")
		if buf.String() != "Tracer package created properly!\n" {
			t.Errorf("Something wrong! Trace should not write '%s'.", buf.String())
		}
	}
}
