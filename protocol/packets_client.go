package protocol

import (
	"io"
	"io/ioutil"
)

//SpawnObject packet
type SpawnObject struct {
	EntityID   VarInt
	ObjectUUID UUID
	Type       byte
	X          int32
	Y          int32
	Z          int32
	Pitch      int8
	Yaw        int8
	Data       int32
	SpeedX     int16
	SpeedY     int16
	SpeedZ     int16
}

func (s *SpawnObject) id() int { return 1 }

func (s *SpawnObject) write(ww io.Writer) (err error) {
	if err = WriteVarInt(ww, s.EntityID); err != nil {
		return
	}
	if err = s.ObjectUUID.Write(ww); err != nil {
		return
	}
	if err = WriteByte(ww, s.Type); err != nil {
		return
	}
	if err = WriteInt32(ww, s.X); err != nil {
		return
	}
	if err = WriteInt32(ww, s.Y); err != nil {
		return
	}
	if err = WriteInt32(ww, s.Z); err != nil {
		return
	}
	if err = WriteInt8(ww, s.Pitch); err != nil {
		return
	}
	if err = WriteInt8(ww, s.Yaw); err != nil {
		return
	}
	if err = WriteInt32(ww, s.Data); err != nil {
		return
	}
	if err = WriteInt16(ww, s.SpeedX); err != nil {
		return
	}
	if err = WriteInt16(ww, s.SpeedY); err != nil {
		return
	}
	if err = WriteInt16(ww, s.SpeedZ); err != nil {
		return
	}
	return
}

func (s *SpawnObject) read(rr io.Reader) (err error) {
	if s.EntityID, err = ReadVarInt(rr); err != nil {
		return
	}
	if err = s.ObjectUUID.Read(rr); err != nil {
		return
	}
	if s.Type, err = ReadByte(rr); err != nil {
		return
	}
	if s.X, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Y, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Z, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Pitch, err = ReadInt8(rr); err != nil {
		return
	}
	if s.Yaw, err = ReadInt8(rr); err != nil {
		return
	}
	if s.Data, err = ReadInt32(rr); err != nil {
		return
	}
	if s.SpeedX, err = ReadInt16(rr); err != nil {
		return
	}
	if s.SpeedY, err = ReadInt16(rr); err != nil {
		return
	}
	if s.SpeedZ, err = ReadInt16(rr); err != nil {
		return
	}
	return
}

//SpawnExperienceOrb packet
type SpawnExperienceOrb struct {
	EntityID VarInt
	X        int32
	Y        int32
	Z        int32
	Count    int16
}

func (s *SpawnExperienceOrb) id() int { return 2 }

func (s *SpawnExperienceOrb) write(ww io.Writer) (err error) {
	if err = WriteVarInt(ww, s.EntityID); err != nil {
		return
	}
	if err = WriteInt32(ww, s.X); err != nil {
		return
	}
	if err = WriteInt32(ww, s.Y); err != nil {
		return
	}
	if err = WriteInt32(ww, s.Z); err != nil {
		return
	}
	if err = WriteInt16(ww, s.Count); err != nil {
		return
	}
	return
}

func (s *SpawnExperienceOrb) read(rr io.Reader) (err error) {
	if s.EntityID, err = ReadVarInt(rr); err != nil {
		return
	}
	if s.X, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Y, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Z, err = ReadInt32(rr); err != nil {
		return
	}
	if s.Count, err = ReadInt16(rr); err != nil {
		return
	}
	return
}

//ServerDifficulty packet
type ServerDifficulty struct {
	Difficulty byte
}

func (s *ServerDifficulty) id() int { return 13 }

func (s *ServerDifficulty) write(ww io.Writer) (err error) {
	if err = WriteByte(ww, s.Difficulty); err != nil {
		return
	}
	return
}

func (s *ServerDifficulty) read(rr io.Reader) (err error) {
	if s.Difficulty, err = ReadByte(rr); err != nil {
		return
	}
	return
}

//PluginMessageClientbound packet
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

//Disconnect packet
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

//JoinGame packet
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
	if err = WriteInt32(ww, j.EntityID); err != nil {
		return
	}
	if err = WriteByte(ww, j.Gamemode); err != nil {
		return
	}
	if err = WriteInt32(ww, j.Dimension); err != nil {
		return
	}
	if err = WriteByte(ww, j.Difficulty); err != nil {
		return
	}
	if err = WriteByte(ww, j.MaxPlayers); err != nil {
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
	if j.EntityID, err = ReadInt32(rr); err != nil {
		return
	}
	if j.Gamemode, err = ReadByte(rr); err != nil {
		return
	}
	if j.Dimension, err = ReadInt32(rr); err != nil {
		return
	}
	if j.Difficulty, err = ReadByte(rr); err != nil {
		return
	}
	if j.MaxPlayers, err = ReadByte(rr); err != nil {
		return
	}
	if j.LevelType, err = ReadString(rr); err != nil {
		return
	}
	if j.ReducedDebugInfo, err = ReadBool(rr); err != nil {
		return
	}
	return
}

type SpawnPosition struct {
	Location Position
}

func (s *SpawnPosition) id() int { return 67 }

func (s *SpawnPosition) write(ww io.Writer) (err error) {
	var tmp [8]byte
	tmp[0] = byte(s.Location >> 56)
	tmp[1] = byte(s.Location >> 48)
	tmp[2] = byte(s.Location >> 40)
	tmp[3] = byte(s.Location >> 32)
	tmp[4] = byte(s.Location >> 24)
	tmp[5] = byte(s.Location >> 16)
	tmp[6] = byte(s.Location >> 8)
	tmp[7] = byte(s.Location >> 0)
	if _, err = ww.Write(tmp[:8]); err != nil {
		return
	}
	return
}

func (s *SpawnPosition) read(rr io.Reader) (err error) {
	var tmp [8]byte
	if _, err = rr.Read(tmp[:8]); err != nil {
		return
	}
	s.Location = (Position(tmp[7]) << 0) | (Position(tmp[6]) << 8) | (Position(tmp[5]) << 16) | (Position(tmp[4]) << 24) | (Position(tmp[3]) << 32) | (Position(tmp[2]) << 40) | (Position(tmp[1]) << 48) | (Position(tmp[0]) << 56)
	return
}

func init() {
	packetList[Play][Clientbound][1] = func() Packet { return &SpawnObject{} }
	packetList[Play][Clientbound][2] = func() Packet { return &SpawnExperienceOrb{} }
	packetList[Play][Clientbound][13] = func() Packet { return &ServerDifficulty{} }
	packetList[Play][Clientbound][24] = func() Packet { return &PluginMessageClientbound{} }
	packetList[Play][Clientbound][26] = func() Packet { return &Disconnect{} }
	packetList[Play][Clientbound][35] = func() Packet { return &JoinGame{} }
	packetList[Play][Clientbound][67] = func() Packet { return &SpawnPosition{} }
}
