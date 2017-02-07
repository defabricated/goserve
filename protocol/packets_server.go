package protocol

import (
	"io"
	"io/ioutil"
)

//PluginMessageServerbound packet
type PluginMessageServerbound struct {
	Channel string
	Data    []byte
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

func init() {
	packetList[Play][Serverbound][9] = func() Packet { return &PluginMessageServerbound{} }
}
