package protocol

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"reflect"
	"sync"

	"github.com/hurisheng/go-futu-api/tcp"
	"google.golang.org/protobuf/proto"
)

type header struct {
	HeaderFlag   [2]byte  // 包头起始标志，固定为“FT”
	ProtoID      uint32   // 协议ID
	ProtoFmtType uint8    // 协议格式类型，0为Protobuf格式，1为Json格式
	ProtoVer     uint8    // 协议版本，用于迭代兼容, 目前填0
	SerialNo     uint32   // 包序列号，用于对应请求包和回包, 要求递增
	BodyLen      uint32   // 包体长度
	BodySHA1     [20]byte // 包体原始数据(解密后)的SHA1哈希值
	Reserved     [8]byte  // 保留8字节扩展
}

type FutuEncoder struct {
	proto  uint32
	serial uint32
	msg    proto.Message
}

var _ tcp.Encoder = (*FutuEncoder)(nil)

func NewEncoder(proto uint32, serial uint32, msg proto.Message) *FutuEncoder {
	return &FutuEncoder{
		proto:  proto,
		serial: serial,
		msg:    msg,
	}
}

func (en *FutuEncoder) WriteTo(c net.Conn) error {
	// 序列化message
	b, err := proto.Marshal(en.msg)
	if err != nil {
		return err
	}
	// 创建header
	h := header{
		HeaderFlag:   [2]byte{'F', 'T'},
		ProtoID:      en.proto,
		ProtoFmtType: 0,
		ProtoVer:     0,
		SerialNo:     en.serial,
		BodyLen:      uint32(len(b)),
	}
	s := sha1.Sum(b)
	for i, c := range s {
		h.BodySHA1[i] = c
	}
	// 将header和body按顺序写入buffer，然后将buffer写入TCP连接
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &h); err != nil {
		return err
	}
	if _, err := buf.Write(b); err != nil {
		return err
	}
	if _, err := buf.WriteTo(c); err != nil {
		return err
	}
	return nil
}

type FutuDecoder struct {
	reg *Registry
}

var _ tcp.Decoder = (*FutuDecoder)(nil)

func NewDecoder(reg *Registry) *FutuDecoder {
	return &FutuDecoder{reg: reg}
}

func (de *FutuDecoder) ReadFrom(c net.Conn) (tcp.Handler, error) {
	// 先读出header，然后根据长度读出body
	// 用registry，proto，serial和body生成当前数据的handler，由tcp框架在gouroutine中处理
	var h header
	// read header
	if err := binary.Read(c, binary.LittleEndian, &h); err != nil {
		return nil, err
	}
	if h.HeaderFlag != [2]byte{'F', 'T'} {
		return nil, errors.New("header flag error")
	}
	// read body
	b := make([]byte, h.BodyLen)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil, err
	}
	// verify body
	s := sha1.Sum(b)
	for i, c := range s {
		if h.BodySHA1[i] != c {
			return nil, errors.New("SHA1 sum error")
		}
	}
	log.Printf("read: proto %v serial %v", h.ProtoID, h.SerialNo)
	return &handler{
		reg:    de.reg,
		proto:  h.ProtoID,
		serial: h.SerialNo,
		body:   b,
	}, nil
}

type handler struct {
	reg    *Registry
	proto  uint32
	serial uint32
	body   []byte
}

var _ tcp.Handler = (*handler)(nil)

func (h *handler) Handle() {
	log.Printf("handle: proto %v serial %v", h.proto, h.serial)
	if err := h.reg.handle(h.proto, h.serial, h.body); err != nil {
		// todo 错误处理
		log.Println(err)
		return
	}
	log.Printf("finish: proto %v serial %v", h.proto, h.serial)
}

type Response interface {
	GetRetType() int32
	GetRetMsg() string
}

func Error(r Response) error {
	if r.GetRetType() != 0 {
		return errors.New(r.GetRetMsg())
	}
	return nil
}

// Registry 接收数据处理器注册表
type Registry struct {
	m  map[uint32]worker
	mu sync.RWMutex
}

// NewRegistry 生成新的Registry
func NewRegistry() *Registry {
	return &Registry{m: make(map[uint32]worker)}
}

// Close 关闭Registry的worker
func (reg *Registry) Close() {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	for _, v := range reg.m {
		v.close()
	}
}

// AddUpdateChan 添加update方法的接收通道
func (reg *Registry) AddUpdateChan(proto uint32, ch RespChan) error {
	return reg.addChan(proto, 0, ch, newUpdateWorker())
}

// AddGetChan 添加get方法的接收通道
func (reg *Registry) AddGetChan(proto uint32, serial uint32, ch RespChan) error {
	return reg.addChan(proto, serial, ch, newGetWorker())
}

func (reg *Registry) addChan(proto uint32, serial uint32, ch RespChan, w worker) error {
	reg.mu.Lock()
	if reg.m[proto] == nil {
		reg.m[proto] = w
	}
	reg.mu.Unlock()
	return reg.m[proto].add(serial, ch)
}

var (
	ErrProtoIDNotFound = errors.New("proto id not found")
)

func (reg *Registry) RemoveChan(proto uint32, serial uint32) error {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	w := reg.m[proto]
	if w == nil {
		return ErrProtoIDNotFound
	}
	return w.remove(serial)
}

func (reg *Registry) handle(proto uint32, serial uint32, body []byte) error {
	reg.mu.RLock()
	defer reg.mu.RUnlock()
	w := reg.m[proto]
	if w == nil {
		return ErrProtoIDNotFound
	}
	return w.handle(serial, body)
}

type RespChan interface {
	Send(b []byte) error
	Close()
}

// 用于接收到数据后，发送协议数据到接收goroutine
type PBChan struct {
	v reflect.Value
	t reflect.Type
}

var _ RespChan = (*PBChan)(nil)

func NewPBChan(i interface{}) (*PBChan, error) {
	// i必须为chan *T类型，T为struct，*T实现proto.Message
	// 通过reflect检查out的类型是否正确
	v, ct := reflect.ValueOf(i), reflect.TypeOf(i)
	// must be a channel type
	if ct.Kind() != reflect.Chan {
		return nil, errors.New("type is not channel")
	}
	// it must be a channel of pointer to the response type which implements proto.Message interface
	pt := ct.Elem()
	if pt.Kind() != reflect.Ptr || !pt.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return nil, errors.New("not a channel of pointer to type implements interface proto.Message")
	}
	return &PBChan{v: v, t: pt.Elem()}, nil
}

func (ch *PBChan) Send(b []byte) error {
	// resp为*T，分配内存空间转换b的数据
	resp := reflect.New(ch.t)
	if err := proto.Unmarshal(b, resp.Interface().(proto.Message)); err != nil {
		return err
	}
	ch.v.Send(resp)
	return nil
}

func (ch *PBChan) Close() {
	ch.v.Close()
}

type worker interface {
	add(serial uint32, ch RespChan) error
	remove(serial uint32) error
	handle(serial uint32, body []byte) error
	close()
}

var (
	ErrDuplicateChannel = errors.New("duplicate channel")
	ErrChannelNotFound  = errors.New("channel not found")
)

// updateWorker 处理update数据推送
type updateWorker struct {
	ch RespChan

	serial uint32
	mu     sync.Mutex
}

var _ worker = (*updateWorker)(nil)

func newUpdateWorker() *updateWorker {
	return &updateWorker{}
}

func (w *updateWorker) add(serial uint32, ch RespChan) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	// channel已经存在，不能重复添加
	if w.ch != nil {
		return ErrDuplicateChannel
	}
	w.serial = 0
	w.ch = ch
	return nil
}

func (w *updateWorker) remove(serial uint32) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ch == nil {
		return ErrChannelNotFound
	}
	w.ch.Close()
	w.serial = 0
	w.ch = nil
	return nil
}

func (w *updateWorker) handle(serial uint32, body []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ch == nil {
		return ErrChannelNotFound
	}
	// serial需递增，已处理过的serial，可能是重复数据
	if w.serial >= serial {
		return errors.New("duplicate serial")
	}
	// 从channel发送，并记录最新的serial
	if err := w.ch.Send(body); err != nil {
		return err
	}
	w.serial = serial
	return nil
}

func (w *updateWorker) close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ch != nil {
		w.ch.Close()
	}
	w.serial = 0
	w.ch = nil
}

type getWorker struct {
	m  map[uint32]RespChan
	mu sync.Mutex
}

var _ worker = (*getWorker)(nil)

func newGetWorker() *getWorker {
	return &getWorker{m: make(map[uint32]RespChan)}
}

func (w *getWorker) add(serial uint32, ch RespChan) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	// serial已经存在，不能重复添加
	if w.m[serial] != nil {
		return ErrDuplicateChannel
	}
	w.m[serial] = ch
	return nil
}

func (w *getWorker) remove(serial uint32) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	ch := w.m[serial]
	if ch == nil {
		return ErrChannelNotFound
	}
	ch.Close()
	delete(w.m, serial)
	return nil
}

func (w *getWorker) handle(serial uint32, body []byte) error {
	// 根据header的serial找到对应的channel，找到返回后，从map中移除
	w.mu.Lock()
	defer w.mu.Unlock()
	// serial不存在，返回错误
	ch := w.m[serial]
	if ch == nil {
		return ErrChannelNotFound
	}
	// 发送数据后，将serial移除
	if err := ch.Send(body); err != nil {
		return err
	}
	ch.Close()
	delete(w.m, serial)
	return nil
}

func (w *getWorker) close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	for _, v := range w.m {
		v.Close()
	}
}
