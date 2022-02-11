package fLogSystem

import (
	"context"
	"sync"
	"time"
)

// 管理者
type SManager struct {
	mutex           sync.RWMutex
	logMutex        sync.Mutex
	cond            *sync.Cond
	ctx             context.Context
	cancel          context.CancelFunc
	serverWG        *sync.WaitGroup
	puberWG         sync.WaitGroup
	publishers      map[int]*SPublisher
	psSerialNo      int
	subscribersReal [LOGLEVELNum][]ISubscriber
	subscribers     [LOGLEVELNum][]ISubscriber
	firstNode       *sPbLinkNode
	lastNode        *sPbLinkNode
}

func NewManager(ctx context.Context, serverWG *sync.WaitGroup) *SManager {
	m := SManager{serverWG: serverWG, publishers: make(map[int]*SPublisher)}
	m.ctx, m.cancel = context.WithCancel(ctx)
	for i := 0; i < int(LOGLEVELNum); i++ {
		m.subscribersReal[i] = make([]ISubscriber, 0)
		m.subscribers[i] = make([]ISubscriber, 0)
	}
	m.cond = sync.NewCond(&m.logMutex)
	m.run()
	return &m
}

func (m *SManager) Init(ctx context.Context, serverWG *sync.WaitGroup) {
	m.serverWG = serverWG
	m.ctx, m.cancel = context.WithCancel(ctx)
	m.publishers = make(map[int]*SPublisher)
	for i := 0; i < int(LOGLEVELNum); i++ {
		m.subscribersReal[i] = make([]ISubscriber, 0)
		m.subscribers[i] = make([]ISubscriber, 0)
	}
	m.cond = sync.NewCond(&m.logMutex)
	m.run()
}

// 創建一個發布者
func (m *SManager) NewPublisher(name string) *SPublisher {
	ps := &SPublisher{name: name, owner: m}
	m.mutex.Lock()
	m.psSerialNo++
	ps.id = m.psSerialNo
	m.publishers[m.psSerialNo] = ps
	m.mutex.Unlock()
	return ps
}

// 訂閱log
func (m *SManager) Subscribe(maskLogLevel int, sb ISubscriber) {
	m.mutex.Lock()
	for i := 0; i < int(LOGLEVELNum); i++ {
		if (maskLogLevel & (1 << i)) != 0 {
			m.subscribers[i] = append(m.subscribers[i], sb)
		}
	}
	m.mutex.Unlock()
}
func (m *SManager) SubscribeReal(maskLogLevel int, sb ISubscriber) {
	m.mutex.Lock()
	for i := 0; i < int(LOGLEVELNum); i++ {
		if (maskLogLevel & (1 << i)) != 0 {
			m.subscribersReal[i] = append(m.subscribersReal[i], sb)
		}
	}
	m.mutex.Unlock()
}

type sPbLinkNode struct {
	log            ILogger
	pb             *SPublisher
	previous, next *sPbLinkNode
}

var pbLogNodePool = sync.Pool{
	New: func() interface{} {
		return new(sPbLinkNode)
	},
}

func (m *SManager) WaitForAllDone() {
	m.puberWG.Wait()
}

func (m *SManager) Cancel() {
	m.cancel()
}

func (m *SManager) Shutdown(timeout int, autoCancel bool) (success bool, cancelled bool) {
	select {
	case <-m.ctx.Done():
		cancelled = true
	default:
		cancelled = false
	}
	success = false
	ch := make(chan struct{}, 1)

	go func() {
		node := pbLogNodePool.Get().(*sPbLinkNode)
		m.logMutex.Lock()
		if m.lastNode != nil {
			node.previous = m.lastNode
			m.lastNode.next = node
			m.lastNode = node
		} else {
			m.firstNode = node
			m.lastNode = node
		}
		m.cond.Signal()
		m.logMutex.Unlock()
		m.WaitForAllDone()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		success = true
		return
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		if autoCancel == false {
			return
		}
		cancelled = true
		m.Cancel()
	}

	if timeout < 100 {
		timeout = 100
	}
	select {
	case <-ch:
		success = true
	case <-time.After(time.Duration(timeout) * time.Millisecond):
	}
	return
}
func (m *SManager) publish(log ILogger, pb *SPublisher) {
	level := log.Level()
	if level >= LOGLEVELNum {
		return
	}
	m.puberWG.Add(1)
	m.mutex.RLock()
	subscribers := m.subscribersReal[level]
	num := len(subscribers)
	for i := 0; i < num; i++ {
		_logPublish(subscribers[i], log, pb, m.ctx)
		select {
		case <-m.ctx.Done():
			m.mutex.RUnlock()
			m.puberWG.Done()
			return
		default:
		}
	}
	m.mutex.RUnlock()

	node := pbLogNodePool.Get().(*sPbLinkNode)
	node.log = log
	node.pb = pb

	m.logMutex.Lock()
	if m.lastNode != nil {
		node.previous = m.lastNode
		m.lastNode.next = node
		m.lastNode = node
	} else {
		m.firstNode = node
		m.lastNode = node
	}
	m.cond.Signal()
	m.logMutex.Unlock()
	m.puberWG.Done()
}

func _logPublish(sb ISubscriber, log ILogger, pb *SPublisher, ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	sb.OnLog(log, pb, ctx)
}

func (m *SManager) run() {
	m.serverWG.Add(1)
	m.puberWG.Add(1)
	go func() {
		var node *sPbLinkNode
		pctx := m.ctx

		for {
			node = nil
			m.cond.L.Lock()
			for m.firstNode == nil {
				// cond.Wait 有個漏洞，不能保證 Wait() 之後不會被別人捷足先登，所以這邊用for而不是if
				// Wait() 會解除 Lock() 等到醒來才重新取得，所以無法保證不會被捷足先登
				m.cond.Wait()
			}
			node = m.firstNode
			m.firstNode = m.firstNode.next
			if m.firstNode != nil {
				m.firstNode.previous = nil
			} else {
				m.lastNode = nil
			}
			node.next = nil
			m.cond.L.Unlock()

			if node.log == nil {
				pbLogNodePool.Put(node)
				node = nil
				m.serverWG.Done()
				m.puberWG.Done()
				return
			}
			select {
			case <-m.ctx.Done():
				pbLogNodePool.Put(node)
				node = nil
				m.serverWG.Done()
				m.puberWG.Done()
				return
			default:
			}
			// 將log發布出去
			// m.pub(node)
			level := node.log.Level()
			if level < LOGLEVELNum {
				m.mutex.RLock()
				subscribers := m.subscribers[level]
				num := len(subscribers)
				for i := 0; i < num; i++ {
					_logPublish(subscribers[i], node.log, node.pb, pctx)
				}
				m.mutex.RUnlock()
			}

			node.log = nil
			node.pb = nil
			pbLogNodePool.Put(node)
		}
	}()
}
