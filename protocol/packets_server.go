package protocol

import (
	"io"
	"io/ioutil"
)

//PluginMessageServerbound packet
type PluginMessageServerbound struct {
	Channel string
	Data    []byte `length:"remaining"`
}

func (p *PluginMessageServerbound) id() int { return 9 }

func (p *PluginMessageServerbound) write(ww io.Writer) (err error) {
	if err = WriteString(ww, p.Channel); err != nil {
		return
	}
	if _, err = ww.Write(p.Data); err != nil {
		return
	}
	return
}

func (p *PluginMessageServerbound) read(rr io.Reader) (err error) {
	if p.Channel, err = ReadString(rr); err != nil {
		return
	}
	if p.Data, err = ioutil.ReadAll(rr); err != nil {
		return
	}
	return
}

//ClientSettings packet
type ClientSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           VarInt
	ChatColors         bool
	DisplayedSkinParts byte
	MainHand           VarInt
}

func (c *ClientSettings) id() int { return 4 }

func (c *ClientSettings) write(ww io.Writer) (err error) {
	if err = WriteString(ww, c.Locale); err != nil {
		return
	}
	if err = WriteByte(ww, c.ViewDistance); err != nil {
		return
	}
	if err = WriteVarInt(ww, c.ChatMode); err != nil {
		return
	}
	if err = WriteBool(ww, c.ChatColors); err != nil {
		return
	}
	if err = WriteByte(ww, c.DisplayedSkinParts); err != nil {
		return
	}
	if err = WriteVarInt(ww, c.MainHand); err != nil {
		return
	}
	return
}

func (c *ClientSettings) read(rr io.Reader) (err error) {
	if c.Locale, err = ReadString(rr); err != nil {
		return
	}
	if c.ViewDistance, err = ReadByte(rr); err != nil {
		return
	}
	if c.ChatMode, err = ReadVarInt(rr); err != nil {
		return
	}
	if c.ChatColors, err = ReadBool(rr); err != nil {
		return
	}
	if c.DisplayedSkinParts, err = ReadByte(rr); err != nil {
		return
	}
	if c.MainHand, err = ReadVarInt(rr); err != nil {
		return
	}
	return
}

func init() {
	packetList[Play][Serverbound][4] = func() Packet { return &ClientSettings{} }
	packetList[Play][Serverbound][9] = func() Packet { return &PluginMessageServerbound{} }
}
