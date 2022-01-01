package fLogSystem

// 發布者介面
type IPublisher interface {
	GetID() int
	GetName() string
}

// 發布者
type SPublisher struct {
	id    int
	name  string
	owner *SManager
}

func (pb *SPublisher) GetID() int {
	return pb.id
}

func (pb *SPublisher) GetName() string {
	return pb.name
}

func (pb *SPublisher) Publish(log ILogger) {
	pb.owner.publish(log, pb)
	// node := pbLogNodePool.Get().(*sPbLinkNode)
	// node.log = log
	// node.pb = pb
	// // 投遞到管理者的queue中
	// owner := pb.owner
	// owner.logMutex.Lock()
	// if owner.lastNode != nil {
	// 	node.previous = owner.lastNode
	// 	owner.lastNode.next = node
	// 	owner.lastNode = node
	// } else {
	// 	owner.firstNode = node
	// 	owner.lastNode = node
	// }
	// owner.cond.Signal()
	// owner.logMutex.Unlock()

	// // level := log.Level()
	// // if level >= LOGLEVELNum {
	// // 	return
	// // }
	// // pb.owner.mutex.Lock()
	// // defer pb.owner.mutex.Unlock()

	// // subscribers := pb.owner.subscribers[level]
	// // num := len(subscribers)
	// // for i := 0; i < num; i++ {
	// // 	subscribers[i].OnLog(log, pb)
	// // }
}
