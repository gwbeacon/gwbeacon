package lib

import (
	"sync/atomic"
	"time"
)

const (
	IndexOffset       = 0
	TimestampOffset   = 19
	ConnectorIDOffset = 50
	IDTypeOffset      = 63
)

const (
	IndexBits       = 19
	TimestampBits   = 31
	ConnectorIDBits = 13
	IDTypeBits      = 1
)

type IDType uint8

const (
	SessionIDType IDType = 0
	MessageIDType IDType = 1
)

type ID uint64

var baseTimestamp uint64 = 0

func SetBaseTimestamp(ts uint64) {
	baseTimestamp = ts
}

type IDMaker interface {
	ConnectorID() uint32
	SetConnectorID(uint32)
	Type() IDType
	MakeID() ID
}

type idMaker struct {
	connID uint32
	idx    uint32
	tp     IDType
}

func NewIDMaker(connID uint32, tp IDType) IDMaker {
	return &idMaker{
		connID: connID,
		tp:     tp,
		idx:    0,
	}
}

func NewSessionIDMaker(connID uint32) IDMaker {
	return NewIDMaker(connID, SessionIDType)
}

func NewMessageIDMaker(connID uint32) IDMaker {
	return NewIDMaker(connID, MessageIDType)
}

func (im *idMaker) Type() IDType {
	return im.tp
}

func (im *idMaker) ConnectorID() uint32 {
	return im.connID
}

func (im *idMaker) SetConnectorID(id uint32) {
	atomic.StoreUint32(&im.connID, id)
}

func (im *idMaker) MakeID() ID {
	idx := atomic.AddUint32((*uint32)(&im.idx), 1)
	var id uint64 = 0
	tm := uint64(time.Now().Unix()) - baseTimestamp
	id |= GetBits(uint64(idx), IndexBits, 0) << IndexOffset
	id |= GetBits(uint64(tm), TimestampBits, 0) << TimestampOffset
	id |= GetBits(uint64(im.connID), ConnectorIDBits, 0) << ConnectorIDOffset
	id |= GetBits(uint64(im.tp), IDTypeBits, 0) << IDTypeOffset
	return ID(id)
}

func Mask(bits uint64, offset uint64) uint64 {
	if bits+offset >= 64 {
		return ^((uint64(1) << offset) - 1)
	}
	return ((uint64(1) << bits) - 1) << offset
}

func GetBits(val uint64, bits uint64, offset uint64) uint64 {
	mask := Mask(bits, offset)
	return (val & mask) >> offset
}

func SetBits(oldVal uint64, bits uint64, offset uint64, val uint64) uint64 {
	return oldVal | GetBits(uint64(val), bits, 0)<<offset
}

func (id *ID) setBits(bits uint64, offset uint64, val uint64) {
	*id = ID(SetBits(uint64(*id), bits, offset, val))
}

func (id ID) getBits(bits uint64, offset uint64) uint64 {
	return GetBits(uint64(id), bits, offset)
}

func (id ID) GetIndex() uint16 {
	return uint16(id.getBits(ConnectorIDBits, ConnectorIDOffset))
}

func (id ID) GetTimestamp() uint64 {
	return uint64(id.getBits(TimestampBits, TimestampOffset))
}

func (id ID) GetConnectorID() uint16 {
	return uint16(id.getBits(ConnectorIDBits, ConnectorIDOffset))
}

func (id ID) IsMessageID() bool {
	return id.getBits(IDTypeBits, IDTypeOffset) == uint64(MessageIDType)
}

func (id ID) IsSessionID() bool {
	return id.getBits(IDTypeBits, IDTypeOffset) == uint64(SessionIDType)
}
