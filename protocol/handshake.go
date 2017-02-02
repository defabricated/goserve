package protocol

type Handshake struct {
	ProtocolVersion VarInt
	Address         string
	Port            uint16
	State           VarInt
}
