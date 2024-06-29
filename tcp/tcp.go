// Package tcp implements the framework of TCP connection.
package tcp

import (
	"errors"
	"io"
	"log"
	"net"

	"golang.org/x/sync/errgroup"
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
	Handle() error
}

type HandlerFunc func() error

func (f HandlerFunc) Handle() error {
	return f()
}

// Conn 根据Protocol读写TCP数据
type Conn struct {
	c  net.Conn
	de Decoder
	g  errgroup.Group
}

// Dial 通过tcp连接目标地址address
func Dial(address string, de Decoder) (*Conn, error) {
	c, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return newConn(c, de), nil
}

func newConn(c net.Conn, de Decoder) *Conn {
	conn := Conn{
		c:  c,
		de: de,
	}
	// 持续从连接读取数据，在单独的goroutine中处理协议返回的Handler
	conn.g.Go(func() error {
		for {
			h, err := conn.de.DecodeFrom(conn.c)
			if err == nil {
				// 成功读取数据，在单独goroutine中执行处理方法
				conn.g.Go(func() error {
					// 保护goroutine不因为panic影响其他goroutine
					defer func() {
						if err := recover(); err != nil {
							log.Printf("panic in handler: %v", err)
						}
					}()
					return h.Handle()
				})
			} else if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				// 如果连接关闭，停止接收数据，其他为数据错误，可忽略
				return nil
			} else {
				log.Printf("decode error: %v", err)
			}
		}
	})
	return &conn
}

// Close 关闭Conn，然后等待已接收数据处理完成
func (conn *Conn) Close() error {
	if err := conn.c.Close(); err != nil {
		return err
	}
	return conn.g.Wait()
}

func (conn *Conn) Write(en Encoder) error {
	return en.EncodeTo(conn.c)
}
