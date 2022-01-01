package fLogSystem

import "context"

// 訂閱者介面
type ISubscriber interface {
	OnLog(logger ILogger, pb *SPublisher, ctx context.Context)
}
