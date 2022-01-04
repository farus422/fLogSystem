package fLogSystem

import (
	"fmt"
)

type SfLogItem struct {
	ItemKey   string
	ItemValue string
}
type SfLog struct {
	SLogger // 以複合（composition）的方式加入SLogger
	caption string
	items   []SfLogItem
}

func NewLog(level LOGLEVEL, format string, param ...interface{}) *SfLog {
	l := SfLog{}
	l.SLogger.Init(level)
	if format != "" {
		l.caption = fmt.Sprintf(format, param...)
	}
	return &l
}

func NewLogEx(level LOGLEVEL, format string, param ...interface{}) *SfLog {
	l := SfLog{}
	l.SLogger.InitAndGetCallstack(level, 1, "")
	if format != "" {
		l.caption = fmt.Sprintf(format, param...)
	}
	return &l
}

func NewLogPanic(level LOGLEVEL, callerAndIgnore string, format string, param ...interface{}) *SfLog {
	l := SfLog{}
	l.SLogger.InitAndPanicCallstack(level, 1, callerAndIgnore)
	if format != "" {
		l.caption = fmt.Sprintf(format, param...)
	}
	return &l
}

func (l *SfLog) SetCaption(format string, param ...interface{}) *SfLog {
	l.caption = fmt.Sprintf(format, param...)
	return l
}

func (l *SfLog) AddItem(key string, value string) *SfLog {
	if l.items == nil {
		l.items = make([]SfLogItem, 0)
	}
	l.items = append(l.items, SfLogItem{ItemKey: key, ItemValue: value})
	return l
}

func (l *SfLog) AddCallstack(skip int, callerAndIgnore string) *SfLog {
	l.SLogger.callstack.GetCallstack(skip+1, callerAndIgnore)
	return l
}

func (l *SfLog) AddPanicCallstack(skip int, callerAndIgnore string) *SfLog {
	l.SLogger.callstack.GetCallstackWithPanic(skip+1, callerAndIgnore)
	return l
}

func (l *SfLog) Message() string {
	msg := l.caption
	if l.items != nil {
		items := l.items
		msg += fmt.Sprintf(". %s=%s", items[0].ItemKey, items[0].ItemValue)
		for i := 1; i < len(items); i++ {
			msg += fmt.Sprintf(", %s=%s", items[i].ItemKey, items[i].ItemValue)
		}
	}
	return msg
}

func (l *SfLog) GetCaption() string {
	return l.caption
}

func (l *SfLog) GetFunctionName() string {
	return l.callstack.GetFunctionName(0)
}

// 範例：列出全部item
// for i, item := 0, log.GetItem(0); item != nil; i, item = i+1, log.GetItem(i+1) {
// 	fmt.Printf("key:%s, value:%s\n", item.ItemKey, item.ItemValue)
// }
func (l *SfLog) GetItem(index int) *SfLogItem {
	if (l.items == nil) || (index >= len(l.items)) {
		return nil
	}
	return &l.items[index]
}
