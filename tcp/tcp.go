// Package tcp implements the framework of TCP connection.
package tcp

import (
	"errors"
	"io"
	"net"
	"sync"
)

// Handler 由Decoder从接收通道中解析数据后返回，执行数据处理任务
type Handler interface {
	Handle()
}

// Encoder 解析类型为[]byte数据，由Conn写入到发送通道
type Encoder interface {
	WriteTo(c net.Conn) error
}

// Decoder 从TCP连接中读取数据的解析器，并返回Handler，在单独goroutine中执行
type Decoder interface {
	// ReadFrom 从TCPConn中读取数据，返回数据处理Handler，如果连接接收缓冲没有数据，方法将会阻塞
	// 返回error保持与TCPConn.Read保持一致，io.EOF为对方关闭，net.ErrClosed为己方关闭
	ReadFrom(c net.Conn) (Handler, error)
}

// Conn 根据Protocol读写TCP数据
type Conn struct {
	c  net.Conn
	de Decoder
	wg sync.WaitGroup
}

// Dial 连接对方
func Dial(network string, address string, de Decoder) (*Conn, error) {
	c, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return newConn(c, de)
}

// Accept 接受连接
func Accept(ln net.Listener, de Decoder) (*Conn, error) {
	c, err := ln.Accept()
	if err != nil {
		return nil, err
	}
	return newConn(c, de)
}

func newConn(c net.Conn, de Decoder) (*Conn, error) {
	conn := Conn{
		c:  c,
		de: de,
	}
	go conn.recv()
	return &conn, nil
}

// Close 关闭TCPConn，然后等待已接收数据处理完成
func (conn *Conn) Close() error {
	if err := conn.c.Close(); err != nil {
		return err
	}
	conn.wg.Wait()
	return nil
}

// recv 持续从连接读取数据，在单独的goroutine中处理协议返回的Handler
func (conn *Conn) recv() {
	for {
		h, err := conn.de.ReadFrom(conn.c)
		if err != nil {
			// 如果连接关闭，停止接收数据，其他为数据错误，可忽略
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}
		} else {
			// 成功读取数据，执行处理方法
			conn.wg.Add(1)
			go func() {
				defer conn.wg.Done()
				h.Handle()
			}()
		}
	}
}

func (conn *Conn) Send(en Encoder) error {
	return en.WriteTo(conn.c)
}
