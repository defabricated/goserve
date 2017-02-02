package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"../protocol"

	"gopkg.in/yaml.v2"
)

//Server contains loaded worlds and players
type Server struct {
	Host string
	Port int

	listener net.Listener
	running  bool

	Motd string

	playerCount int32
	maxPlayers  int32
}

var GoServer *Server

//CreateServer function creates new instance of Server
func CreateServer(host string, port int) *Server {
	data, _ := ioutil.ReadFile("./settings.yml")

	var settings Settings

	yaml.Unmarshal(data, &settings)

	server := &Server{
		Host:        host,
		Port:        port,
		running:     true,
		Motd:        settings.Motd,
		playerCount: 0,
		maxPlayers:  settings.MaxPlayers,
	}

	return server
}

//Start function starts a Minecraft server
func (server *Server) Start() error {
	address := server.Host + ":" + strconv.Itoa(server.Port)

	listen, err := net.Listen("tcp", address)

	if err != nil {
		return err
	}

	server.listener = listen

	GoServer = server

	for {
		conn, err := listen.Accept()

		if err != nil {
			return err
		}

		go server.HandleConnection(conn)
	}
}

//HandleConnection function handles new incoming connection
func (server *Server) HandleConnection(conn net.Conn) {
	log.Println("Incoming connection from " + conn.RemoteAddr().String())

	minecraftConnection := &protocol.Conn{
		Out:            conn,
		In:             conn,
		Deadliner:      conn,
		ReadDirection:  protocol.Serverbound,
		WriteDirection: protocol.Clientbound,
	}

	packet, err := minecraftConnection.ReadPacket()

	if err != nil {
		return
	}

	handshake, ok := packet.(protocol.Handshake)

	if !ok {
		return
	}

	if handshake.State == 1 {
		minecraftConnection.State = protocol.Status
		packet, err := minecraftConnection.ReadPacket()
		if _, ok := packet.(protocol.StatusGet); !ok || err != nil {
			return
		}

		ping := &protocol.Ping{}
		ping.Version = protocol.PingVersion{
			Name:     "GoServe",
			Protocol: 47,
		}

		ping.Players = protocol.PingPlayers{
			Max:    9999,
			Online: int(server.playerCount),
		}

		ping.Description = server.Motd

		by, err := json.Marshal(ping)
		if err != nil {
			return
		}

		statusResponse := protocol.StatusResponse{
			Data: string(by),
		}

		minecraftConnection.WritePacket(statusResponse)

		packet, err = minecraftConnection.ReadPacket()
		if err != nil {
			return
		}
		cPing, ok := packet.(protocol.ClientStatusPing)

		if !ok {
			return
		}
		minecraftConnection.WritePacket(protocol.StatusPing{Time: cPing.Time})
		return
	}
}
