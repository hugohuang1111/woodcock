package gate

import "sync"

var (
	connMap       map[uint64]connect
	connIDMutex   sync.Mutex
	connIDCounter uint64
)

func init() {
	connMap = make(map[uint64]connect)
}

func clientConnect(c connect) {
	connIDMutex.Lock()
	connIDCounter++
	c.setID(connIDCounter)
	connMap[connIDCounter] = c
	connIDMutex.Unlock()
}

func clientDisconnect(c connect) {
	delete(connMap, c.ID())
}
