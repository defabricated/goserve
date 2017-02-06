package protocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"
)

//Conn allows to receive and send Minecraft packets
type Conn struct {
	In  io.Reader
	Out io.Writer

	Deadliner Deadliner

	State                State
	compressionThreshold int

	ReadDirection  Direction
	WriteDirection Direction
}

type Deadliner interface {
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}

//ReadPacket function reads a minecraft packet from connection
func (conn *Conn) ReadPacket() (Packet, error) {
	if conn.Deadliner != nil {
		conn.Deadliner.SetReadDeadline(time.Now().Add(10 * time.Second))
	}

	size, err := ReadVarInt(conn.In)
	if err != nil {
		return nil, err
	}

	if size < 0 {
		return nil, errors.New("Invalid length: negative")
	}

	buf := make([]byte, size)

	if _, err = io.ReadFull(conn.In, buf); err != nil {
		return nil, err
	}

	var reader *bytes.Reader
	reader = bytes.NewReader(buf)

	//TODO Packet compression

	id, err := ReadVarInt(reader)

	if err != nil {
		return nil, err
	}

	packets := packetList[conn.State][conn.ReadDirection]

	if id < 0 || int(id) >= len(packets) || packets[id] == nil {
		return nil, fmt.Errorf("Unknown packet %s:%02X", "", id)
	}

	packet := packets[id]()

	if err := packet.read(reader); err != nil {
		return packet, fmt.Errorf("packet(%s:%02X): %s", "", id, err)
	}

	if reader.Len() > 0 {
		return packet, fmt.Errorf("Cannot finish reading packet %s:%02X, have %d bytes left", "", id, reader.Len())
	}

	return packet, nil
}

//WritePacket function writes the packet to connection
func (conn *Conn) WritePacket(packet Packet) error {
	if conn.Deadliner != nil {
		conn.Deadliner.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}

	buf := &bytes.Buffer{}

	if err := WriteVarInt(buf, VarInt(packet.id())); err != nil {
		return err
	}

	if err := packet.write(buf); err != nil {
		return err
	}

	if err := WriteVarInt(conn.Out, VarInt(buf.Len())); err != nil {
		return err
	}

	_, err := buf.WriteTo(conn.Out)
	return err
}
