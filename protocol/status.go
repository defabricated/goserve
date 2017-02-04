package protocol

import (
	"encoding/json"
	"io"
)

type StatusResponse struct {
	Data *Ping `as:"json"`
}

func (s *StatusResponse) id() int { return 0 }

func (s *StatusResponse) write(ww io.Writer) (err error) {
	var tmp0 []byte
	if tmp0, err = json.Marshal(&s.Data); err != nil {
		return
	}
	tmp1 := string(tmp0)
	if err = WriteString(ww, tmp1); err != nil {
		return
	}

	return
}

func (s *StatusResponse) read(rr io.Reader) (err error) {
	var tmp0 string
	if tmp0, err = ReadString(rr); err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(tmp0), &s.Data); err != nil {
		return
	}
	return
}

type StatusPing struct {
	Time int64
}

func (s *StatusPing) id() int { return 1 }

func (s *StatusPing) write(ww io.Writer) (err error) {
	var tmp [8]byte
	tmp[0] = byte(s.Time >> 56)
	tmp[1] = byte(s.Time >> 48)
	tmp[2] = byte(s.Time >> 40)
	tmp[3] = byte(s.Time >> 32)
	tmp[4] = byte(s.Time >> 24)
	tmp[5] = byte(s.Time >> 16)
	tmp[6] = byte(s.Time >> 8)
	tmp[7] = byte(s.Time >> 0)
	if _, err = ww.Write(tmp[:8]); err != nil {
		return
	}
	return
}

func (s *StatusPing) read(rr io.Reader) (err error) {
	var tmp [8]byte
	if _, err = rr.Read(tmp[:8]); err != nil {
		return
	}
	s.Time = int64((uint64(tmp[7]) << 0) | (uint64(tmp[6]) << 8) | (uint64(tmp[5]) << 16) | (uint64(tmp[4]) << 24) | (uint64(tmp[3]) << 32) | (uint64(tmp[2]) << 40) | (uint64(tmp[1]) << 48) | (uint64(tmp[0]) << 56))
	return
}

type StatusGet struct {
}

func (s *StatusGet) id() int { return 0 }
func (s *StatusGet) write(ww io.Writer) (err error) {
	return
}
func (s *StatusGet) read(rr io.Reader) (err error) {
	return
}

type ClientStatusPing struct {
	Time int64
}

func (c *ClientStatusPing) id() int { return 1 }

func (c *ClientStatusPing) write(ww io.Writer) (err error) {
	var tmp [8]byte
	tmp[0] = byte(c.Time >> 56)
	tmp[1] = byte(c.Time >> 48)
	tmp[2] = byte(c.Time >> 40)
	tmp[3] = byte(c.Time >> 32)
	tmp[4] = byte(c.Time >> 24)
	tmp[5] = byte(c.Time >> 16)
	tmp[6] = byte(c.Time >> 8)
	tmp[7] = byte(c.Time >> 0)
	if _, err = ww.Write(tmp[:8]); err != nil {
		return
	}
	return
}

func (c *ClientStatusPing) read(rr io.Reader) (err error) {
	var tmp [8]byte
	if _, err = rr.Read(tmp[:8]); err != nil {
		return
	}
	c.Time = int64((uint64(tmp[7]) << 0) | (uint64(tmp[6]) << 8) | (uint64(tmp[5]) << 16) | (uint64(tmp[4]) << 24) | (uint64(tmp[3]) << 32) | (uint64(tmp[2]) << 40) | (uint64(tmp[1]) << 48) | (uint64(tmp[0]) << 56))
	return
}

func init() {
	packetList[Status][Clientbound][0] = func() Packet { return &StatusResponse{} }
	packetList[Status][Clientbound][1] = func() Packet { return &StatusPing{} }

	packetList[Status][Serverbound][0] = func() Packet { return &StatusGet{} }
	packetList[Status][Serverbound][1] = func() Packet { return &ClientStatusPing{} }
}
