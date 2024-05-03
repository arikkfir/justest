package internal

import "runtime"

//go:noinline
func runtimeGetStack(skip int) []runtime.Frame {
	var pcs [32]uintptr
	frames, n := callers(skip, pcs[:])
	ff := make([]runtime.Frame, 0, n)
	for {
		fr, ok := frames.Next()
		if !ok {
			break
		}
		ff = append(ff, fr)
	}
	return ff
}

//go:noinline
func runtimeGetFrame(skip int) runtime.Frame {
	var pcs [3]uintptr
	frames, _ := callers(skip, pcs[:])
	fr, ok := frames.Next()
	if !ok {
		return runtime.Frame{}
	}
	return fr
}

//go:noinline
func callers(skip int, pcs []uintptr) (frames *runtime.Frames, n int) {
	n = runtime.Callers(skip+1, pcs)
	frames = runtime.CallersFrames(pcs[:n])
	if _, ok := frames.Next(); !ok {
		return &runtime.Frames{}, 0
	}
	return
}
