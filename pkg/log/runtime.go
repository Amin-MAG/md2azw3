package ravandlog

import (
	"context"
	"fmt"
	"runtime"
	"strings"
)

const LoggerPackageName = "github.com/Amin-MAG/md2azw3/pkg/log"
const MaximumNumberOfStackFrames = 5

func GetCurrentFunctionName(ctx context.Context) context.Context {
	skip := 0
	for skip <= MaximumNumberOfStackFrames {
		// Get the program counter, file name, and line number
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			Warn("cannot determine function name from stack trace")
			return ctx
		}

		// Get the function object for this caller
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			return ctx
		}

		// Check if the function is inside the log package
		if strings.Contains(fn.Name(), LoggerPackageName) {
			skip++
			continue
		}

		// Add the File name and line number
		ctx = context.WithValue(ctx, ContextKeyFile, fmt.Sprintf("%s:%d", file, line))
		// Add the Function name
		ctx = context.WithValue(ctx, ContextKeyFunction, fn.Name())

		return ctx
	}

	return ctx
}
