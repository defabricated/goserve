package protocol

import "io"

type Handshake struct {
	ProtocolVersion VarInt
	Address         string
	Port            uint16
	State           VarInt
}

func (handshake *Handshake) id() int {
	return 0
}

func (handshake *Handshake) write(writer io.Writer) (err error) {
	var tmp [2]byte
	if err = WriteVarInt(writer, handshake.ProtocolVersion); err != nil {
		return
	}
	if err = WriteString(writer, handshake.Address); err != nil {
		return
	}
	tmp[0] = byte(handshake.Port >> 8)
	tmp[1] = byte(handshake.Port >> 0)
	if _, err = writer.Write(tmp[:2]); err != nil {
		return
	}
	if err = WriteVarInt(writer, handshake.State); err != nil {
		return
	}

	return
}

func (handshake *Handshake) read(reader io.Reader) (err error) {
	var tmp [2]byte
	if handshake.ProtocolVersion, err = ReadVarInt(reader); err != nil {
		return
	}
	if handshake.Address, err = ReadString(reader); err != nil {
		return
	}
	if _, err = reader.Read(tmp[:2]); err != nil {
		return
	}
	handshake.Port = ((uint16(tmp[1]) << 0) | (uint16(tmp[0]) << 8))
	if handshake.State, err = ReadVarInt(reader); err != nil {
		return
	}
	return
}

func init() {
	packetList[Handshaking][Serverbound][0] = func() Packet { return &Handshake{} }
}
