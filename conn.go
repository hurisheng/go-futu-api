package futuapi

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

type protoHeader struct {
	HeaderFlag   [2]byte  // 包头起始标志，固定为“FT”
	ProtoID      uint32   // 协议ID
	ProtoFmtType uint8    // 协议格式类型，0为Protobuf格式，1为Json格式
	ProtoVer     uint8    // 协议版本，用于迭代兼容, 目前填0
	SerialNo     uint32   // 包序列号，用于对应请求包和回包, 要求递增
	BodyLen      uint32   // 包体长度
	BodySHA1     [20]byte // 包体原始数据(解密后)的SHA1哈希值
	Reserved     [8]byte  // 保留8字节扩展
}

// Predefined errors.
var (
	ErrConnClosed = fmt.Errorf("connection is closed")
)

// Conn is socket connection to server daemon. It manages socket communication with a list of receiving operation.
// There are 2 kinds of receiving, one time and continuously.
type conn struct {
	socket net.Conn

	queue      map[uint32]respHandler
	lastSerial uint32
	lock       sync.Mutex

	closed bool
	wg     sync.WaitGroup
}

// NewConn connects to server daemon.
func newConn(address string) (*conn, error) {
	socket, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("connect socket error: %w", err)
	}

	c := &conn{
		socket: socket,
		queue:  make(map[uint32]respHandler),
	}
	go c.receive()
	return c, nil
}

// Close disconnects from server and close all pending channel.
func (c *conn) close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.closed = true
	if err := c.socket.Close(); err != nil {
		return fmt.Errorf("close socket error: %w", err)
	}

	c.wg.Wait()
	for _, h := range c.queue {
		h.close()
	}

	return nil
}

func (c *conn) receive() {
	for {
		header := new(protoHeader)
		// read header
		if err := binary.Read(c.socket, binary.LittleEndian, header); err != nil {
			break
		}
		if header.HeaderFlag == [2]byte{'F', 'T'} {
			// read body
			body := make([]byte, header.BodyLen)
			if _, err := io.ReadFull(c.socket, body); err != nil {
				break
			}
			c.wg.Add(1)
			go func() {
				defer c.wg.Done()
				// verify body
				s := sha1.Sum(body)
				for i, c := range s {
					if header.BodySHA1[i] != c {
						return
					}
				}
				// pass to handler, if protoID is not registered, just ignore
				h := c.queue[header.ProtoID]
				if h != nil {
					h.handle(header.SerialNo, body)
				}
			}()
		}
	}
}

// Send sends data to server and return the receiving buffer immediately. When data is returned from server, data is returned via channel.
// Send can be call multiple times for same protoID.
func (c *conn) send(protoID uint32, body []byte) (<-chan []byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return nil, ErrConnClosed
	}

	// create header
	header := &protoHeader{
		HeaderFlag:   [2]byte{'F', 'T'},
		ProtoID:      protoID,
		ProtoFmtType: 0,
		ProtoVer:     0,
		SerialNo:     c.lastSerial + 1,
		BodyLen:      uint32(len(body)),
	}
	s := sha1.Sum(body)
	for i, c := range s {
		header.BodySHA1[i] = c
	}
	// write header and body to buffer
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, header); err != nil {
		return nil, fmt.Errorf("marshall header error: %w", err)
	}
	if n, err := buf.Write(body); err != nil || n != len(body) {
		return nil, fmt.Errorf("write body to buffer error: %w, %d", err, n)
	}
	// create handler if not exists
	if c.queue[protoID] == nil {
		c.queue[protoID] = newMsgRespHandler()
	}
	// add receiving channel to queue of protoID, ready for receiving
	h, ok := c.queue[protoID].(*msgRespHandler)
	if !ok {
		return nil, fmt.Errorf("handler type error")
	}
	out := make(chan []byte)
	h.add(c.lastSerial+1, out)
	// write to socket
	if _, err := buf.WriteTo(c.socket); err != nil {
		h.remove(c.lastSerial + 1)
		return nil, fmt.Errorf("send to socket error: %w", err)
	}
	c.lastSerial++
	return out, nil
}

// Subscribe returns the receiving channel for protoID. This function can be call only once for same protoID.
// When incoming data arrived, a new buffer will be sent to the channel.
func (c *conn) subscribe(protoID uint32) (<-chan []byte, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return nil, ErrConnClosed
	}

	if c.queue[protoID] != nil {
		return nil, fmt.Errorf("subscribe to duplicated protoID")
	}
	out := make(chan []byte)
	c.queue[protoID] = newNotifyRespHandler(out)
	return out, nil
}
