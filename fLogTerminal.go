package fLogSystem

import (
	"context"
	"fmt"

	color "github.com/fatih/color"
)

type STerminal struct{}

func NewTerminal() *STerminal {
	return &STerminal{}
}

var logLevelToColor = [LOGLEVELNum]color.Attribute{color.FgHiBlack, color.FgWhite, color.FgWhite, color.FgGreen, color.FgMagenta, color.FgRed, color.FgRed}

func (sb *STerminal) OnLog(logger ILogger, pb *SPublisher, ctx context.Context) {
	ltime := logger.Time()
	level := logger.Level()
	if color.NoColor {
		fmt.Printf("\033[%dm%d-%02d-%02d %02d:%02d:%02d.%03d | %s | %s:%d %s() | \033[%dm%s\x1b[0m\n", color.FgBlue, ltime.Year(), ltime.Month(), ltime.Day(), ltime.Hour(), ltime.Minute(), ltime.Second(), ltime.Nanosecond()/1000000, LOGTypeNameShot[logger.Level()], logger.Filename(), logger.Line(), logger.FunctionName(), [logger.Level()], logger.Message())
	} else {
		cbase := color.New(color.FgHiBlue)
		cbase.Printf("%d-%02d-%02d %02d:%02d:%02d.%03d | %s | %s:%d %s() | ", ltime.Year(), ltime.Month(), ltime.Day(), ltime.Hour(), ltime.Minute(), ltime.Second(), ltime.Nanosecond()/1000000, LOGTypeNameShot[logger.Level()], logger.Filename(), logger.Line(), logger.FunctionName())
		cMsg := color.New(logLevelToColor[level])
		cMsg.Printf("%s\n", logger.Message())
	}

	// fmt.Printf("%d-%d-%d %d:%d:%d | %s | \033[%d;%d;%dm%s\x1b[0m\n", ltime.Year(), ltime.Month(), ltime.Day(), ltime.Hour(), ltime.Minute(), ltime.Second(), flog.LOGTypeNameShot[logger.Level()], 0, color.BgYellow, logLevelToColor[logger.Level()], logger.Message())
	// fmt.Printf("%d-%d-%d %d:%d:%d | %s | \033[%dm%s\x1b[0m\n", ltime.Year(), ltime.Month(), ltime.Day(), ltime.Hour(), ltime.Minute(), ltime.Second(), flog.LOGTypeNameShot[logger.Level()], logLevelToColor[logger.Level()], logger.Message())
	// fmt.Println(logger.Message())
	switch logger.(type) {
	case *SfLog:
		// log := logger.(*SfLog)
		// for i, item := 0, log.GetItem(0); item != nil; i, item = i+1, log.GetItem(i+1) {
		// 	fmt.Printf("key:%s, value:%s\n", item.ItemKey, item.ItemValue)
		// }
		callers := logger.Callstack()
		if callers != nil {
			callNum := len(callers)
			if color.NoColor {
				fmt.Printf("\x1b[%dm           呼叫堆疊如下：\n", color.FgCyan)
				for i := 0; i < callNum; i++ {
					fmt.Printf("\x1b[%dm           %s:%d \x1b[%dm%s()\x1b[0m\n", color.FgYellow, callers[i].File, callers[i].Line, color.FgCyan, callers[i].Function)
				}
			} else {
				cCyan := color.New(color.FgCyan)
				cCyan.Println("           呼叫堆疊如下：")
				cCaller := color.New(color.FgYellow)
				for i := 0; i < callNum; i++ {
					cCaller.Printf("           %s:%d ", callers[i].File, callers[i].Line)
					cCyan.Printf("%s()\n", callers[i].Function)
				}
			}
		}
	}
	// if strings.Compare("伺服器關機", logger.Message()) == 0 {
	// 	time.Sleep(time.Second * 30)
	// }
}
