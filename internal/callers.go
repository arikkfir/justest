package internal

// CallerAt returns a Frame that describes a frame on the caller's
// stack. The argument skipCaller is the number of frames to skip over.
//
//go:noinline
func CallerAt(skipCallers int) Frame {
	return getFrame(skipCallers + 3)
}

// CallStackAt returns all the Frames that describe the caller's stack.
// The argument skipCaller is the number of frames to skip over.
//
//go:noinline
func CallStackAt(skipCallers int) Frames {
	st := getStack(skipCallers + 3)
	ff := make(Frames, len(st))
	for i, fr := range st {
		ff[i] = fr
	}
	return ff
}
