package internal

import "runtime"

// getFrame translates a runtime.Frame item returned from the internal
// runtime utilities into a frame.
//
//go:noinline
func getFrame(skipCallers int) *frame {
	return &frame{pc: runtimeGetFrame(skipCallers).PC}
}

// Frame defines an interface for accessing and displaying stack frame
// information for debugging, optimizing or inspection. Usually you will
// find Frame in a Frames slice, acting as a stack trace or stack dump.
//
// Frames are meant to be seen, so we have implemented the following
// default formatting verbs on it:
//
//	"%s"  – the base name of the file (or `unknown`) and the line number (if known)
//	"%q"  – the same as `%s` but wrapped in `"` delimiters
//	"%d"  – the line number
//	"%n"  – the basic function name, ie without a full package qualifier
//	"%v"  – the full path of the file (or `unknown`) and the line number (if known)
//	"%+v" – a standard line in a stack trace: a full function name on one line,
//	        and a full file name and line number on a second line
//	"%#v" – a Golang representation with the type (`errors.Frame`)
//
// Marshaling a frame as text uses the `%+v` format.
// Marshaling as JSON returns an object with location data:
//
//	{"function":"test.pkg.in/example.init","file":"/src/example.go","line":10}
//
// A Frame is immutable, so no setters are provided, but you can copy
// one trivially with:
//
//	function, file, line := oldFrame.Location()
//	newFrame := errors.NewFrame(function, file, line)
type Frame interface {
	// Location returns the frame's caller's characteristics for help with
	// identifying and debugging the codebase.
	//
	// Location results are generated uniquely per Frame implementation.
	// When using this package's implementation, note that the results are
	// evaluated and expanded lazily when the frame was generated from the
	// local call stack: Location is not safe for concurrent access.
	Location() (function string, file string, line int)
}

// frame is this package's default implementation of Frame in such a way
// that we can create one either from the actual call stack or
// "synthetically:" by parsing a stack trace or even specifically
// designating the location characteristics. frame also implements
// interfaces to integrate with runtime (via program counters) and
// serialization and deserialization processes.
type frame struct {
	pc        uintptr
	runtimeFn *runtime.Func
	function  string
	file      string
	line      int
}

// Location returns the frame's caller's characteristics for help with
// identifying and debugging the codebase.
//
// The results are evaluated and expanded lazily when the frame was
// generated from the local call stack: Location is not safe for
// concurrent access.
//
//go:noinline
func (f *frame) Location() (function string, file string, line int) {
	return f.getFunction(), f.getFile(), f.getLine()
}

// getFunction gets the frame's full caller function name. Prioritizes
// synthetic values if available, otherwise expands the pc using runtime
// and memorizes the result.
//
//go:noinline
func (f *frame) getFunction() (function string) {
	function = f.function
	if function == "" {
		function = "unknown"
		if f.pc != 0 {
			function = f.fn().Name()
			f.function = function
		}
	}
	return
}

// getFile gets the frame's caller's filename. Prioritizes synthetic
// values if available, otherwise expands the pc using runtime and
// memorizes the result.
//
//go:noinline
func (f *frame) getFile() (file string) {
	file = f.file
	if file == "" {
		file = "unknown"
		if f.pc != 0 {
			file, _ = f.fn().FileLine(f.pc)
			f.file = file
		}
	}
	return
}

// getLine gets the frame's caller's file line. Prioritizes synthetic
// values if available, otherwise expands the pc using runtime and
// memorizes the result.
//
//go:noinline
func (f *frame) getLine() (line int) {
	line = f.line
	if line == 0 {
		if f.pc != 0 {
			_, line = f.fn().FileLine(f.pc)
			f.line = line
		}
	}
	return
}

// fn is the way to cleanly access the runtimeFn field: if none is found
// it attempts to look it up from the frame location program counter
// (pc). This lookup will only happen once.
//
//go:noinline
func (f *frame) fn() *runtime.Func {
	if f.runtimeFn == nil && f.pc != 0 {
		f.runtimeFn = runtime.FuncForPC(f.pc)
	}
	return f.runtimeFn
}
