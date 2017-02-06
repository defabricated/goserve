package entity

import (
	"time"

	"../message"
	"../protocol"
)

type Player struct {
	conn *protocol.Conn

	uuid string
	Name string

	packetQueue   chan protocol.Packet
	packetRead    chan protocol.Packet
	errorChannel  chan error
	ClosedChannel chan struct{}

	ping int32
}

func NewPlayer(name string, uuid string, conn *protocol.Conn) *Player {
	player := &Player{
		conn:        conn,
		packetQueue: make(chan protocol.Packet, 200),
		packetRead:  make(chan protocol.Packet, 20),
		uuid:        uuid,
		Name:        name,
		ping:        -1,
	}

	go player.packetReader()
	go player.packetWriter()

	return player
}

func (player *Player) Join() {
	player.QueuePacket(&protocol.JoinGame{
		EntityID:   int32(0),
		Gamemode:   byte(0),
		Dimension:  int32(0),
		Difficulty: byte(0),
		MaxPlayers: byte(0),
		LevelType:  "default",
	})

	player.QueuePacket(&protocol.PluginMessageClientbound{
		Channel: "MC|Brand",
		Data:    []byte("GoServe"),
	})

	time.Sleep(time.Second * 3)

	player.QueuePacket(&protocol.Disconnect{
		Data: (&message.Message{Text: "Successfully joined to server!", Color: message.Blue}).JSONString(),
	})
}

func (player *Player) QueuePacket(packet protocol.Packet) {
	select {
	case player.packetQueue <- packet:
	case <-player.ClosedChannel:
	}
}

func (player *Player) packetReader() {
	for {
		packet, err := player.conn.ReadPacket()
		if err != nil {
			player.errorChannel <- err
			return
		}
		select {
		case player.packetRead <- packet:
		case <-player.ClosedChannel:
			return
		}
	}
}

func (player *Player) packetWriter() {
	for {
		select {
		case packet := <-player.packetQueue:
			player.conn.WritePacket(packet)
		case <-player.ClosedChannel:
			return
		}
	}
}
