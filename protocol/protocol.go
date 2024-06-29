// protocol 为消息添加header信息，从tcp连接中解析数据，转换protobuf数据和[]byte格式
package protocol

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/hurisheng/go-futu-api/tcp"
	"google.golang.org/protobuf/proto"
)

type FutuProtocol struct {
	conn *tcp.Conn
	de   *decoder

	serial atomic.Uint32 //当前的序列号
}

func Connect(address string, api map[uint32]Worker) (*FutuProtocol, error) {
	de := newDecoder(api)
	conn, err := tcp.Dial(address, de)
	if err != nil {
		return nil, err
	}
	return &FutuProtocol{
		conn: conn,
		de:   de,
	}, nil
}

func (ft *FutuProtocol) SerialNo() uint32 {
	return ft.serial.Add(1)
}

func (ft *FutuProtocol) RegisterGet(proto uint32, req proto.Message, out *ProtobufChan) error {
	se := ft.SerialNo()
	if err := ft.de.register(proto, se, out); err != nil {
		return err
	}
	if err := ft.conn.Write(newEncoder(proto, se, req)); err != nil {
		if err := ft.de.register(proto, se, nil); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (ft *FutuProtocol) RegisterUpdate(proto uint32, out *ProtobufChan) error {
	return ft.de.register(proto, 0, out)
}

func (ft *FutuProtocol) Close() error {
	ft.de.Close()
	return ft.conn.Close()
}

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

type encoder struct {
	proto  uint32
	serial uint32
	msg    proto.Message
}

func newEncoder(proto uint32, serial uint32, msg proto.Message) *encoder {
	return &encoder{
		proto:  proto,
		serial: serial,
		msg:    msg,
	}
}

func (en *encoder) EncodeTo(w io.Writer) error {
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
		BodySHA1:     sha1.Sum(b),
	}
	// 将header和body按顺序写入buffer，然后将buffer写入io连接
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.LittleEndian, &h); err != nil {
		return err
	}
	if _, err := buf.Write(b); err != nil {
		return err
	}
	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

type decoder struct {
	reg map[uint32]Worker

	mu sync.RWMutex
}

func newDecoder(api map[uint32]Worker) *decoder {
	return &decoder{
		reg: api,
	}
}

func (de *decoder) DecodeFrom(r io.Reader) (tcp.Handler, error) {
	// 先读出header，然后根据长度读出body
	// 用registry，proto，serial和body生成当前数据的handler，由tcp框架在gouroutine中处理
	var h header
	// read header
	if err := binary.Read(r, binary.LittleEndian, &h); err != nil {
		return nil, err
	}
	if h.HeaderFlag != [2]byte{'F', 'T'} {
		return nil, errors.New("header flag error")
	}
	// read body
	b := make([]byte, h.BodyLen)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	// verify body
	if h.BodySHA1 != sha1.Sum(b) {
		return nil, errors.New("SHA1 sum error")
	}
	//log.Printf("read: proto %v serial %v", h.ProtoID, h.SerialNo)
	return tcp.HandlerFunc(func() error {
		//log.Printf("handle: proto %v serial %v", h.ProtoID, h.SerialNo)
		de.mu.RLock()
		defer de.mu.RUnlock()
		if de.reg[h.ProtoID] == nil {
			return errors.New("worker is not found")
		}
		return de.reg[h.ProtoID].handle(h.SerialNo, b)
	}), nil
}

func (de *decoder) Close() {
	de.mu.RLock()
	defer de.mu.RUnlock()
	for k := range de.reg {
		de.reg[k].close()
	}
}

func (de *decoder) register(proto uint32, serial uint32, out *ProtobufChan) error {
	de.mu.RLock()
	defer de.mu.RUnlock()
	if de.reg[proto] == nil {
		return errors.New("worker is not found")
	}
	return de.reg[proto].register(serial, out)
}

type Worker interface {
	register(serial uint32, out *ProtobufChan) error
	handle(serial uint32, body []byte) error
	close()
}

// updater 处理update数据并推送到指定的chan
type updater struct {
	ch *ProtobufChan // 发送消息的chan

	serial uint32 // 记录最后接收到的序列号，防止重复数据，新数据序列号递增
	mu     sync.RWMutex
}

func NewUpdater() Worker {
	return &updater{}
}

func (w *updater) register(_ uint32, out *ProtobufChan) error {
	// 如果已经存在chan，关闭，然后替换为新的out
	if w.ch != nil {
		w.ch.close()
	}
	w.ch = out
	return nil
}

func (w *updater) handle(serial uint32, body []byte) error {
	// serial需递增，已处理过的serial，可能是重复数据
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.serial >= serial {
		return errors.New("might be outdated message received")
	}
	if w.ch == nil {
		return errors.New("sending out channel is nil")
	}
	// 解析消息，并从chan发送，更新serial为最新接收的
	if err := w.ch.send(body); err != nil {
		return err
	}
	w.serial = serial
	return nil
}

func (w *updater) close() {
	w.ch.close()
}

type getter struct {
	m map[uint32]*getterItem

	mu sync.RWMutex
}

func NewGetter() Worker {
	return &getter{
		m: make(map[uint32]*getterItem),
	}
}

func (w *getter) register(serial uint32, out *ProtobufChan) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	// serial已经存在，关闭channel，替换为新的out
	if w.m[serial] != nil {
		w.m[serial].close()
	}
	w.m[serial] = newGetterItem(out)
	return nil
}

func (w *getter) handle(serial uint32, body []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	// 根据header的serial找到对应的getterItem，处理完后，从map中移除
	if w.m[serial] == nil {
		return errors.New("getter item does not exist")
	}
	// 解析消息，并从chan发送，删除getterItem
	if err := w.m[serial].handle(body); err != nil {
		return err
	}
	delete(w.m, serial)
	return nil
}

func (w *getter) close() {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for k := range w.m {
		w.m[k].close()
	}
}

type getterItem struct {
	ch *ProtobufChan
}

func newGetterItem(out *ProtobufChan) *getterItem {
	return &getterItem{
		ch: out,
	}
}

func (i *getterItem) handle(body []byte) error {
	if i.ch == nil {
		return errors.New("sending out channel is nil")
	}
	if err := i.ch.send(body); err != nil {
		return err
	}
	i.ch.close()
	return nil
}

func (i *getterItem) close() {
	i.ch.close()
}

// 用于接收到数据后，发送协议数据到接收channel
type ProtobufChan struct {
	v reflect.Value
	t reflect.Type
}

// 从chan *T类型转换为ProtobufChan类型，T为struct，*T实现proto.Message
func NewProtobufChan(i any) *ProtobufChan {
	// i必须为chan *T类型，T为struct，*T实现proto.Message
	// 通过reflect检查类型是否正确
	ct := reflect.TypeOf(i)
	// must be a channel type
	if ct.Kind() != reflect.Chan {
		return nil
	}
	// it must be a channel of pointer to the response type which implements proto.Message interface
	pt := ct.Elem()
	if pt.Kind() != reflect.Ptr || !pt.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		return nil
	}
	st := pt.Elem()
	if st.Kind() != reflect.Struct {
		return nil
	}
	return &ProtobufChan{v: reflect.ValueOf(i), t: st}
}

func (ch *ProtobufChan) send(b []byte) error {
	// resp为*T，分配内存空间转换b的数据
	buf := reflect.New(ch.t)
	if err := proto.Unmarshal(b, buf.Interface().(proto.Message)); err != nil {
		return err
	}
	ch.v.Send(buf)
	return nil
}

func (ch *ProtobufChan) close() {
	if ch != nil {
		ch.v.Close()
	}
}

// Response是protobuf接口定义的返回信息获取方法
type Response interface {
	GetRetType() int32
	GetRetMsg() string
}

// Error将Response转换为error类型
func Error(r Response) error {
	if r.GetRetType() != 0 {
		return errors.New(r.GetRetMsg())
	}
	return nil
}
