package protocol

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
