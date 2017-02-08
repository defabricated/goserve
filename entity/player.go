package entity

import (
	"errors"
	"log"
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

	brand string

	settings struct {
		locale string
	}
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

	player.QueuePacket(&protocol.ServerDifficulty{
		Difficulty: 0,
	})

	player.QueuePacket(&protocol.SpawnPosition{
		Location: protocol.NewPosition(0, 64, 0),
	})

	player.QueuePacket(&protocol.PlayerAbilities{
		Flags:        0,
		FlyingSpeed:  1,
		WalkingSpeed: 1,
	})

	tick := time.NewTicker(time.Second / 10)
	defer tick.Stop()

	for {
		select {
		case err := <-player.errorChannel:
			log.Printf("Player %s error: %s\n", player.Name, err)
			return
		case packet := <-player.packetRead:
			player.handlePacket(packet)
		}
	}
}

//Disconnect function disconnects the player from the server
func (player *Player) Disconnect(reason message.Message) {
	player.QueuePacket(&protocol.Disconnect{reason})
	player.errorChannel <- errors.New(reason.Text)
}

//QueuePacket queues a packet to be sent to the player
func (player *Player) QueuePacket(packet protocol.Packet) {
	select {
	case player.packetQueue <- packet:
	case <-player.ClosedChannel:
	}
}

func (player *Player) handlePacket(packet protocol.Packet) {
	switch packet := packet.(type) {
	case *protocol.ClientSettings:
		player.settings.locale = packet.Locale

		player.Disconnect(message.Message{Text: "Successfully joined to server!", Color: message.Blue})
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
