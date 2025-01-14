package berr

import (
	"fmt"
	"runtime"
	"strings"
)

func caller(skip int) *frame {
	pc, _, _, _ := runtime.Caller(skip) // callers(3) skips this method, stack.callers, and runtime.Callers
	f := frame(pc)
	return &f
}

// frame is a single program counter of a stack frame.
type frame uintptr

func (f frame) pc() uintptr {
	return uintptr(f)
}

func (f frame) get() StackFrame {
	pc := f.pc()
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()

	i := strings.LastIndex(frame.Function, "/")
	name := frame.Function[i+1:] // e.g pkg.(*receiver).MethodName

	return StackFrame{
		Name: name,
		File: frame.File,
		Line: frame.Line,
	}
}

type Stack []StackFrame

func (s Stack) Compact(opt ...Option) []string {
	var result []string
	for _, v := range s {
		result = append(result, v.Print(opt...))
	}
	return result
}

type StackFrame struct {
	Name string
	File string
	Line int
}

func (sf StackFrame) String() string {
	return fmt.Sprintf("%s:%d:%s", sf.File, sf.Line, sf.Name)
}

func (sf StackFrame) Print(opts ...Option) string {
	opt := defaultOptions()
	for _, o := range opts {
		o(opt)
	}
	file, line, name := "", "", sf.Name
	if opt.showFile {
		if opt.shortFile {
			file = ": " + sf.ShortFile()
		} else {
			file = ": " + sf.File
		}
	}
	if opt.showLine {
		line = fmt.Sprintf(":%d", sf.Line)
	}
	if opt.shortName {
		name = sf.ShortFuncName()
	}
	return fmt.Sprintf("%s%s%s", name, file, line)
}

func (sf StackFrame) ShortFuncName() string {
	token := strings.Split(sf.Name, ".")
	return token[len(token)-1]
}

func (sf StackFrame) ShortFile() string {
	token := strings.Split(sf.File, "/")
	return token[len(token)-1]
}

func getStackFrame(skip int) *StackFrame {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()
	frame := caller(skip)
	stackFrame := frame.get()
	return &stackFrame
}
