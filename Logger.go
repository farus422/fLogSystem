package fLogSystem

import (
	"path"
	"runtime"
	"strings"
	"time"

	fcb "github.com/farus422/fCallstack"
)

// Logger
type ILogger interface {
	Level() LOGLEVEL
	Time() time.Time
	Message() string
	Filename() string
	Line() int
	FunctionName() string
	Callstack() []fcb.SCaller
}

type SLogger struct {
	level     LOGLEVEL
	time      time.Time
	caller    fcb.SCaller
	callstack fcb.SCallstack
}

func (l *SLogger) Init(level LOGLEVEL, skip int) {
	l.time = time.Now()
	l.level = level
	// l.callstack.Clean()
	if pc, file, line, ok := runtime.Caller(skip + 1); ok {
		_, l.caller.File = path.Split(file)
		l.caller.Line = line
		function := runtime.FuncForPC(pc)
		if function != nil {
			// l.caller.Function = function.Name()
			_, funcname := path.Split(function.Name())
			funcs := strings.SplitAfterN(funcname, ".", 2)
			l.caller.Function = funcs[1]
		}
	} else {
		l.caller.Line = 0
		l.caller.File = "--< fetch failed >--"
		l.caller.Function = "--< fetch failed >--"
	}
}

func (l *SLogger) InitAndGetCallstack(level LOGLEVEL, skip int, callerAndIgnore string) {
	l.time = time.Now()
	l.level = level
	// l.callstack.Clean()
	l.callstack.GetCallstack(skip+1, callerAndIgnore)
	callers := l.callstack.GetCallers()
	if callers != nil {
		_, l.caller.File = path.Split(callers[0].File)
		l.caller.Line = callers[0].Line
		l.caller.Function = callers[0].Function
	}
}

func (l *SLogger) InitAndPanicCallstack(level LOGLEVEL, skip int, callerAndIgnore string) {
	l.time = time.Now()
	l.level = level
	// l.callstack.Clean()
	l.callstack.GetCallstackWithPanic(skip+1, callerAndIgnore)
	callers := l.callstack.GetCallers()
	if callers != nil {
		_, l.caller.File = path.Split(callers[0].File)
		l.caller.Line = callers[0].Line
		l.caller.Function = callers[0].Function
	}
}

// func (l *SLogger) InitAndGetCallstack(level LOGLEVEL, skip int, callerAndIgnore string) {
// 	l.time = time.Now()
// 	l.level = level
// 	l.callstack.GetCallstack(skip+1, callerAndIgnore)
// }

// func (l *SLogger) InitWithPanic(level LOGLEVEL, skip int, callerAndIgnore string) {
// 	l.time = time.Now()
// 	l.level = level
// 	l.callstack.GetCallstackWithPanic(skip+1, callerAndIgnore)
// }

func (l *SLogger) Level() LOGLEVEL {
	return l.level
}

func (l *SLogger) Time() time.Time {
	return l.time
}

func (l *SLogger) Filename() string {
	return l.caller.File
}

func (l *SLogger) Line() int {
	return l.caller.Line
}

func (l *SLogger) FunctionName() string {
	return l.caller.Function
	// callers := l.callstack.GetCallers()
	// if callers == nil {
	// 	return ""
	// }
	// return callers[0].Function
}

func (l *SLogger) Callstack() []fcb.SCaller {
	return l.callstack.GetCallers()
}

func (l *SLogger) Message() string {
	return ""
}
