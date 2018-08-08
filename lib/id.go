package lib

import "time"

const (
	IndexOffset       = 0
	PortOffset        = 0
	ReservedOffset    = 16
	TimestampOffset   = 19
	ConnectorIDOffset = 50
	IDTypeOffset      = 63
)

const (
	IndexBits       = 19
	PortBits        = 16
	ReservedBits    = 3
	TimestampBits   = 31
	ConnectorIDBits = 13
	IDTypeBits      = 1
)

const (
	SessionIDType = 0
	MessageIDType = 1
)

type ID uint64

func MakeMessageID(connID uint16, idx uint32) ID {
	var id uint64 = 0
	tm := uint64(time.Now().Unix())
	id |= GetBits(uint64(idx), IndexBits, 0) << IndexOffset
	id |= GetBits(uint64(tm), TimestampBits, 0) << TimestampOffset
	id |= GetBits(uint64(connID), ConnectorIDBits, 0) << ConnectorIDOffset
	id |= GetBits(uint64(MessageIDType), IDTypeBits, 0) << IDTypeOffset
	return ID(id)
}

func MakeSessionID(connID uint16, port uint16) ID {
	var id uint64 = 0
	tm := uint64(time.Now().Unix())
	id |= GetBits(uint64(port), PortBits, 0) << PortOffset
	id |= GetBits(uint64(tm), TimestampBits, 0) << TimestampOffset
	id |= GetBits(uint64(connID), ConnectorIDBits, 0) << ConnectorIDOffset
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

func (id *ID) SetBits(bits uint64, offset uint64, val uint64) {
	*id = ID(SetBits(uint64(*id), bits, offset, val))
}

func (id ID) GetBits(bits uint64, offset uint64) uint64 {
	return GetBits(uint64(id), bits, offset)
}

func (id ID) GetIndex() uint16 {
	return uint16(id.GetBits(ConnectorIDBits, ConnectorIDOffset))
}

func (id ID) GetPort() uint16 {
	return uint16(id.GetBits(PortBits, PortOffset))
}

func (id ID) GetReserved() uint16 {
	return uint16(id.GetBits(ReservedBits, ReservedOffset))
}

func (id ID) GetConnectorID() uint16 {
	return uint16(id.GetBits(ConnectorIDBits, ConnectorIDOffset))
}

func (id ID) IsMessageID() bool {
	return id.GetBits(IDTypeBits, IDTypeOffset) == MessageIDType
}

func (id ID) IsSessionID() bool {
	return id.GetBits(IDTypeBits, IDTypeOffset) == SessionIDType
}
