package lib

import "net"

type ServerType string

const (
	RegisterServer  ServerType = "register"
	ConnectorServer            = "connector"
	UserServer                 = "user"
	ChatServer                 = "chat"
	RosterServer               = "roster"
	MUCServer                  = "muc"
)

type Server interface {
	Type() ServerType
	SetID(id uint16)
	Serve(lis net.Listener) error
}
