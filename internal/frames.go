package internal

// getStack translates runtime.Frame items returned from the internal
// runtime utilities into frames.
//
//go:noinline
func getStack(skipCallers int) frames {
	st := runtimeGetStack(skipCallers)
	ff := make([]*frame, len(st))
	for i, fr := range st {
		ff[i] = &frame{pc: fr.PC}
	}
	return ff
}

// Frames is a slice of Frame data. This can represent a stack trace or
// some subset of a stack trace.
type Frames []Frame

// frames stores a slice of frame structs and implements both the
// StackFrames and stackTracer interfaces.
type frames []*frame
