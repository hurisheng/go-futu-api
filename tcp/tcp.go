// Package tcp implements the framework of TCP connection.
package tcp

import (
	"errors"
	"io"
	"net"
	"sync"
)

// Encoder 是tcp连接的编码器
type Encoder interface {
	// EncodeTo ，将数据编码为字节流，写入到tcp连接
	EncodeTo(w io.Writer) error
}

// Decoder 是tcp连接的解码器
type Decoder interface {
	// DecodeFrom 从tcp连接中读取字节流数据进行解码，返回数据处理器Handler，数据处理器在单独goroutine中运行。
	// 如果tcp连接接收缓冲没有数据，方法将会阻塞。
	// 返回 error 保持与net包一致，io.EOF为对方关闭，net.ErrClosed为己方关闭
	DecodeFrom(r io.Reader) (Handler, error)
}

// Handler 处理解码后的数据
type Handler interface {
	Handle()
}

// Conn 根据Protocol读写TCP数据
type Conn struct {
	c  net.Conn
	de Decoder
	wg sync.WaitGroup
}

// Dial 通过tcp连接目标地址address
func Dial(address string, de Decoder) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return newConn(c, de)
}

// Listen 在地址address上监听tcp连接
func Listen(address string) (net.Listener, error) {
	return net.Listen("tcp", address)
}

// Accept 等待并返回新的tcp连接
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

// Close 关闭Conn，然后等待已接收数据处理完成
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
		h, err := conn.de.DecodeFrom(conn.c)
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

func (conn *Conn) Write(en Encoder) error {
	return en.EncodeTo(conn.c)
}
