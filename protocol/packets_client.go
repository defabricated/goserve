package protocol

import (
	"io"
	"io/ioutil"
)

type PluginMessageClientbound struct {
	Channel string
	Data    []byte
}

func (p *PluginMessageClientbound) id() int { return 24 }

func (p *PluginMessageClientbound) write(ww io.Writer) (err error) {
	if err = WriteString(ww, p.Channel); err != nil {
		return
	}
	if _, err = ww.Write(p.Data); err != nil {
		return
	}
	return
}

func (p *PluginMessageClientbound) read(rr io.Reader) (err error) {
	if p.Channel, err = ReadString(rr); err != nil {
		return
	}
	if p.Data, err = ioutil.ReadAll(rr); err != nil {
		return
	}
	return
}

type Disconnect struct {
	Data string
}

func (d *Disconnect) id() int { return 26 }

func (d *Disconnect) write(ww io.Writer) (err error) {
	if err = WriteString(ww, d.Data); err != nil {
		return
	}
	return
}

func (d *Disconnect) read(rr io.Reader) (err error) {
	if d.Data, err = ReadString(rr); err != nil {
		return
	}
	return
}

type JoinGame struct {
	EntityID         int32
	Gamemode         byte
	Dimension        int32
	Difficulty       byte
	MaxPlayers       byte
	LevelType        string
	ReducedDebugInfo bool
}

func (j *JoinGame) id() int { return 35 }

func (j *JoinGame) write(ww io.Writer) (err error) {
	var tmp [4]byte
	tmp[0] = byte(j.EntityID >> 24)
	tmp[1] = byte(j.EntityID >> 16)
	tmp[2] = byte(j.EntityID >> 8)
	tmp[3] = byte(j.EntityID >> 0)
	if _, err = ww.Write(tmp[:4]); err != nil {
		return
	}
	tmp[0] = byte(j.Gamemode >> 0)
	if _, err = ww.Write(tmp[:1]); err != nil {
		return
	}
	tmp[0] = byte(j.Dimension >> 24)
	tmp[1] = byte(j.Dimension >> 16)
	tmp[2] = byte(j.Dimension >> 8)
	tmp[3] = byte(j.Dimension >> 0)
	if _, err = ww.Write(tmp[:4]); err != nil {
		return
	}
	tmp[0] = byte(j.Difficulty >> 0)
	if _, err = ww.Write(tmp[:1]); err != nil {
		return
	}
	tmp[0] = byte(j.MaxPlayers >> 0)
	if _, err = ww.Write(tmp[:1]); err != nil {
		return
	}
	if err = WriteString(ww, j.LevelType); err != nil {
		return
	}
	if err = WriteBool(ww, j.ReducedDebugInfo); err != nil {
		return
	}
	return
}

func (j *JoinGame) read(rr io.Reader) (err error) {
	var tmp [4]byte
	if _, err = rr.Read(tmp[:4]); err != nil {
		return
	}
	j.EntityID = int32((uint32(tmp[3]) << 0) | (uint32(tmp[2]) << 8) | (uint32(tmp[1]) << 16) | (uint32(tmp[0]) << 24))
	if _, err = rr.Read(tmp[:1]); err != nil {
		return
	}
	j.Gamemode = (byte(tmp[0]) << 0)
	if _, err = rr.Read(tmp[:1]); err != nil {
		return
	}
	j.Dimension = int32((uint32(tmp[3]) << 0) | (uint32(tmp[2]) << 8) | (uint32(tmp[1]) << 16) | (uint32(tmp[0]) << 24))
	if _, err = rr.Read(tmp[:1]); err != nil {
		return
	}
	j.Difficulty = (byte(tmp[0]) << 0)
	if _, err = rr.Read(tmp[:1]); err != nil {
		return
	}
	j.MaxPlayers = (byte(tmp[0]) << 0)
	if j.LevelType, err = ReadString(rr); err != nil {
		return
	}
	if j.ReducedDebugInfo, err = ReadBool(rr); err != nil {
		return
	}
	return
}

func init() {
	packetList[Play][Clientbound][24] = func() Packet { return &PluginMessageClientbound{} }
	packetList[Play][Clientbound][26] = func() Packet { return &Disconnect{} }
	packetList[Play][Clientbound][35] = func() Packet { return &JoinGame{} }
}
