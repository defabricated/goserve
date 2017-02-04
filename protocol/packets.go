package protocol

import (
	"io"
)

type State int

const (
	Handshaking State = iota
	Play
	Login
	Status
)

type Direction int

const (
	Serverbound Direction = iota
	Clientbound
)

type Packet interface {
	id() int
	write(io.Writer) error
	read(io.Reader) error
}

var packetList [4][2][100]func() Packet
