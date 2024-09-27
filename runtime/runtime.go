package runtime

import (
	"runtime"
)

func GetFrame(skip int) runtime.Frame {
	pc := make([]uintptr, 10)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]
	frame, _ := runtime.CallersFrames(pc).Next()

	return runtime.Frame{
		Function: frame.Function,
		File:     frame.File,
		Line:     frame.Line,
	}
}
