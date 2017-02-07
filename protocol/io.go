package protocol

import (
	"errors"
	"io"
	"math"
)

type VarInt int32
type VarLong int64
type UUID [16]byte

const varPart = uint32(0x7F)
const varPartLong = uint64(0x7F)

type Position uint64

func NewPosition(x, y, z int) Position {
	return ((Position(x) & 0x3FFFFFF) << 38) | ((Position(y) & 0xFFF) << 26) | (Position(z) & 0x3FFFFFF)
}

func (p Position) X() int {
	return int(int64(p) >> 38)
}

func (p Position) Y() int {
	return int((int64(p) >> 26) & 0xFFF)
}

func (p Position) Z() int {
	return int(int64(p) << 38 >> 38)
}

var (
	ErrVarIntTooLarge  = errors.New("VarInt too large")
	ErrVarLongTooLarge = errors.New("VarLong too large")
)

func varIntSize(i VarInt) int {
	size := 0
	ui := uint32(i)
	for {
		size++
		if (ui & ^varPart) == 0 {
			return size
		}
		ui >>= 7
	}

}

func WriteVarInt(w io.Writer, i VarInt) error {
	ui := uint32(i)
	for {
		if (ui & ^varPart) == 0 {
			err := WriteByte(w, byte(ui))
			return err
		}
		err := WriteByte(w, byte((ui&varPart)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

func ReadVarInt(r io.Reader) (VarInt, error) {
	var size uint
	var val uint32
	for {
		b, err := ReadByte(r)
		if err != nil {
			return VarInt(val), err
		}

		val |= (uint32(b) & varPart) << (size * 7)
		size++
		if size > 5 {
			return VarInt(val), ErrVarIntTooLarge
		}

		if (b & 0x80) == 0 {
			break
		}
	}
	return VarInt(val), nil
}

func WriteInt32(w io.Writer, i int32) error {
	var tmp [4]byte
	tmp[0] = byte(i >> 24)
	tmp[1] = byte(i >> 16)
	tmp[2] = byte(i >> 8)
	tmp[3] = byte(i >> 0)
	if _, err := w.Write(tmp[:4]); err != nil {
		return err
	}
	return nil
}

func ReadInt32(r io.Reader) (int32, error) {
	var tmp [4]byte
	if _, err := r.Read(tmp[:4]); err != nil {
		return 0, err
	}
	i := int32((uint32(tmp[3]) << 0) | (uint32(tmp[2]) << 8) | (uint32(tmp[1]) << 16) | (uint32(tmp[0]) << 24))
	return i, nil
}

func WriteInt16(w io.Writer, i int16) error {
	var tmp [2]byte
	tmp[0] = byte(i >> 8)
	tmp[1] = byte(i >> 0)
	if _, err := w.Write(tmp[:2]); err != nil {
		return err
	}
	return nil
}

func ReadInt16(r io.Reader) (int16, error) {
	var tmp [2]byte
	if _, err := r.Read(tmp[:2]); err != nil {
		return 0, err
	}
	i := int16((uint16(tmp[1]) << 0) | (uint16(tmp[0]) << 8))
	return i, nil
}

func WriteInt8(w io.Writer, i int8) error {
	err := WriteByte(w, byte(i))
	return err
}

func ReadInt8(r io.Reader) (int8, error) {
	b, err := ReadByte(r)
	return int8(b), err
}

func WriteVarLong(w io.Writer, i VarLong) error {
	ui := uint64(i)
	for {
		if (ui & ^varPartLong) == 0 {
			err := WriteByte(w, byte(ui))
			return err
		}
		err := WriteByte(w, byte((ui&varPartLong)|0x80))
		if err != nil {
			return err
		}
		ui >>= 7
	}
}

func ReadVarLong(r io.Reader) (VarLong, error) {
	var size uint
	var val uint64
	for {
		b, err := ReadByte(r)
		if err != nil {
			return VarLong(val), err
		}

		val |= (uint64(b) & varPartLong) << (size * 7)
		size++
		if size > 10 {
			return VarLong(val), ErrVarLongTooLarge
		}

		if (b & 0x80) == 0 {
			break
		}
	}
	return VarLong(val), nil
}

func WriteString(w io.Writer, str string) error {
	b := []byte(str)
	err := WriteVarInt(w, VarInt(len(b)))
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func ReadString(r io.Reader) (string, error) {
	l, err := ReadVarInt(r)
	if err != nil {
		return "", nil
	}
	if l < 0 || l > math.MaxInt16 {
		return "", errors.New("string length out of bounds")
	}
	buf := make([]byte, int(l))
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func WriteBool(w io.Writer, b bool) error {
	if b {
		return WriteByte(w, 1)
	}
	return WriteByte(w, 0)
}

func ReadBool(r io.Reader) (bool, error) {
	b, err := ReadByte(r)
	if b == 0 {
		return false, err
	}
	return true, err
}

func WriteByte(w io.Writer, b byte) error {
	if bw, ok := w.(io.ByteWriter); ok {
		return bw.WriteByte(b)
	}
	var buf [1]byte
	buf[0] = b
	_, err := w.Write(buf[:1])
	return err
}

func ReadByte(r io.Reader) (byte, error) {
	if br, ok := r.(io.ByteReader); ok {
		return br.ReadByte()
	}
	var buf [1]byte
	_, err := r.Read(buf[:1])
	return buf[0], err
}

func (u *UUID) Write(w io.Writer) error {
	_, err := w.Write(u[:])
	return err
}

func (u *UUID) Read(r io.Reader) error {
	_, err := io.ReadFull(r, u[:])
	return err
}
