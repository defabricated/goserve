package protocol

import (
	"fmt"
	"io"
	"math"
)

type LoginDisconnect struct {
	Data string
}

func (loginDisconnect *LoginDisconnect) id() int {
	return 0
}

func (loginDisconnect *LoginDisconnect) write(writer io.Writer) (err error) {
	if err = WriteString(writer, loginDisconnect.Data); err != nil {
		return
	}
	return
}

func (loginDisconnect *LoginDisconnect) read(reader io.Reader) (err error) {
	if loginDisconnect.Data, err = ReadString(reader); err != nil {
		return
	}
	return
}

type EncryptionKeyRequest struct {
	ServerID    string
	PublicKey   []byte
	VerifyToken []byte
}

func (e *EncryptionKeyRequest) id() int {
	return 1
}

func (e *EncryptionKeyRequest) write(writer io.Writer) (err error) {
	if err = WriteString(writer, e.ServerID); err != nil {
		return
	}
	if err = WriteVarInt(writer, VarInt(len(e.PublicKey))); err != nil {
		return
	}
	if _, err = writer.Write(e.PublicKey); err != nil {
		return
	}
	if err = WriteVarInt(writer, VarInt(len(e.VerifyToken))); err != nil {
		return
	}
	if _, err = writer.Write(e.VerifyToken); err != nil {
		return
	}
	return
}

func (e *EncryptionKeyRequest) read(reader io.Reader) (err error) {
	if e.ServerID, err = ReadString(reader); err != nil {
		return
	}
	var tmp0 VarInt
	if tmp0, err = ReadVarInt(reader); err != nil {
		return
	}
	if tmp0 > math.MaxInt16 {
		return fmt.Errorf("array larger than max value: %d > %d", tmp0, math.MaxInt16)
	}
	if tmp0 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp0)
	}
	e.PublicKey = make([]byte, tmp0)
	if _, err = reader.Read(e.PublicKey); err != nil {
		return
	}
	var tmp1 VarInt
	if tmp1, err = ReadVarInt(reader); err != nil {
		return
	}
	if tmp1 > math.MaxInt16 {
		return fmt.Errorf("array larger than max value: %d > %d", tmp1, math.MaxInt16)
	}
	if tmp1 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp1)
	}
	e.VerifyToken = make([]byte, tmp1)
	if _, err = reader.Read(e.VerifyToken); err != nil {
		return
	}
	return
}

type LoginSuccess struct {
	UUID     string
	Username string
}

func (l *LoginSuccess) id() int { return 2 }

func (l *LoginSuccess) write(ww io.Writer) (err error) {
	if err = WriteString(ww, l.UUID); err != nil {
		return
	}
	if err = WriteString(ww, l.Username); err != nil {
		return
	}
	return
}

func (l *LoginSuccess) read(rr io.Reader) (err error) {
	if l.UUID, err = ReadString(rr); err != nil {
		return
	}
	if l.Username, err = ReadString(rr); err != nil {
		return
	}
	return
}

type LoginStart struct {
	Username string
}

func (l *LoginStart) id() int { return 0 }

func (l *LoginStart) write(ww io.Writer) (err error) {
	if err = WriteString(ww, l.Username); err != nil {
		return
	}
	return
}
func (l *LoginStart) read(rr io.Reader) (err error) {
	if l.Username, err = ReadString(rr); err != nil {
		return
	}
	return
}

type EncryptionKeyResponse struct {
	SharedSecret []byte
	VerifyToken  []byte
}

func (e *EncryptionKeyResponse) id() int { return 1 }

func (e *EncryptionKeyResponse) write(writer io.Writer) (err error) {
	if err = WriteVarInt(writer, VarInt(len(e.SharedSecret))); err != nil {
		return
	}
	if _, err = writer.Write(e.SharedSecret); err != nil {
		return
	}
	if err = WriteVarInt(writer, VarInt(len(e.VerifyToken))); err != nil {
		return
	}
	if _, err = writer.Write(e.VerifyToken); err != nil {
		return
	}
	return
}

func (e *EncryptionKeyResponse) read(reader io.Reader) (err error) {
	var tmp0 VarInt
	if tmp0, err = ReadVarInt(reader); err != nil {
		return
	}
	if tmp0 > math.MaxInt16 {
		return fmt.Errorf("array larger than max value: %d > %d", tmp0, math.MaxInt16)
	}
	if tmp0 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp0)
	}
	e.SharedSecret = make([]byte, tmp0)
	if _, err = reader.Read(e.SharedSecret); err != nil {
		return
	}
	var tmp1 VarInt
	if tmp1, err = ReadVarInt(reader); err != nil {
		return
	}
	if tmp1 > math.MaxInt16 {
		return fmt.Errorf("array larger than max value: %d > %d", tmp1, math.MaxInt16)
	}
	if tmp1 < 0 {
		return fmt.Errorf("negative array size: %d < 0", tmp1)
	}
	e.VerifyToken = make([]byte, tmp1)
	if _, err = reader.Read(e.VerifyToken); err != nil {
		return
	}
	return
}

func init() {
	packetList[Login][Clientbound][0] = func() Packet { return &LoginDisconnect{} }
	packetList[Login][Clientbound][1] = func() Packet { return &EncryptionKeyRequest{} }
	packetList[Login][Clientbound][2] = func() Packet { return &LoginSuccess{} }

	packetList[Login][Serverbound][0] = func() Packet { return &LoginStart{} }
	packetList[Login][Serverbound][1] = func() Packet { return &EncryptionKeyResponse{} }
}
