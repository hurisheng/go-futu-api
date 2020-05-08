package futuapi

import (
	"sync"
)

type respHandler interface {
	handle(serialNo uint32, body []byte)
	close()
}

type notifyRespHandler struct {
	out        chan []byte
	lastSerial uint32
	lock       sync.Mutex
	closed     bool
}

func newNotifyRespHandler(out chan []byte) *notifyRespHandler {
	return &notifyRespHandler{
		out: out,
	}
}

func (h *notifyRespHandler) handle(serialNo uint32, body []byte) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return
	}
	// serialNo less than the last one, just ignore
	if serialNo <= h.lastSerial {
		return
	}
	h.lastSerial = serialNo
	h.out <- body
}

func (h *notifyRespHandler) close() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.closed = true
	close(h.out)
}

type msgRespHandler struct {
	queue  map[uint32]chan []byte // serialNo is the key for direct access.
	lock   sync.Mutex
	closed bool
}

func newMsgRespHandler() *msgRespHandler {
	return &msgRespHandler{
		queue: make(map[uint32]chan []byte),
	}
}

func (h *msgRespHandler) handle(serialNo uint32, body []byte) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return
	}
	if h.queue[serialNo] == nil {
		return
	}
	h.queue[serialNo] <- body
	delete(h.queue, serialNo)
}

func (h *msgRespHandler) close() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.closed = true
	for _, out := range h.queue {
		close(out)
	}
}

func (h *msgRespHandler) add(serialNo uint32, out chan []byte) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return
	}
	if h.queue[serialNo] != nil {
		return
	}
	h.queue[serialNo] = out
}

func (h *msgRespHandler) remove(serialNo uint32) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if h.closed {
		return
	}
	delete(h.queue, serialNo)
}
