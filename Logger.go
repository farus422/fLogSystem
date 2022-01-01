package fLogSystem

import (
	"time"

	fcb "github.com/farus422/fCallstack"
)

// Logger
type ILogger interface {
	Level() LOGLEVEL
	Time() time.Time
	Message() string
	FunctionName() string
	Callstack() []fcb.SCaller
}

type SLogger struct {
	level     LOGLEVEL
	time      time.Time
	callstack fcb.SCallstack
}

func (l *SLogger) Init(level LOGLEVEL) {
	l.time = time.Now()
	l.level = level
	l.callstack.Clean()
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

func (l *SLogger) FunctionName() string {
	callers := l.callstack.GetCallers()
	if callers == nil {
		return ""
	}
	return callers[0].Function
}
func (l *SLogger) Callstack() []fcb.SCaller {
	return l.callstack.GetCallers()
}

func (l *SLogger) Message() string {
	return ""
}
