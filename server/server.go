package server

import (
	"io/ioutil"
	"log"
	"net"
	"strconv"

	"../entity/"
	"../message/"
	"../protocol"
	"../protocol/auth"

	"gopkg.in/yaml.v2"
)

//Server contains loaded worlds and players
type Server struct {
	Host string
	Port int

	listener net.Listener
	running  bool

	protocolVersion int

	authenticator auth.Authenticator

	Motd string

	playerCount int32
	maxPlayers  int32

	onlineMode bool
}

var GoServer *Server

//CreateServer function creates new instance of Server
func CreateServer(host string, port int) *Server {
	data, _ := ioutil.ReadFile("./settings.yml")

	var settings Settings

	yaml.Unmarshal(data, &settings)

	server := &Server{
		Host:            host,
		Port:            port,
		running:         true,
		protocolVersion: 316,
		authenticator:   auth.Instance,
		Motd:            settings.Motd,
		playerCount:     0,
		maxPlayers:      settings.MaxPlayers,
		onlineMode:      settings.OnlineMode,
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

	handshake, ok := packet.(*protocol.Handshake)

	if !ok {
		return
	}

	if handshake.State == 1 {
		minecraftConnection.State = protocol.Status
		packet, _ := minecraftConnection.ReadPacket()

		if _, ok := packet.(*protocol.StatusGet); !ok || err != nil {
			return
		}

		ping := &protocol.Ping{}
		ping.Version = protocol.PingVersion{
			Name:     "GoServe",
			Protocol: server.protocolVersion,
		}

		ping.Players = protocol.PingPlayers{
			Max:    int(server.maxPlayers),
			Online: int(server.playerCount),
		}

		ping.Description = server.Motd

		statusResponse := &protocol.StatusResponse{
			Data: ping,
		}

		minecraftConnection.WritePacket(statusResponse)

		packet, err = minecraftConnection.ReadPacket()
		if err != nil {
			return
		}

		cPing, ok := packet.(*protocol.ClientStatusPing)

		if !ok {
			return
		}
		minecraftConnection.WritePacket(&protocol.StatusPing{Time: cPing.Time})
		return
	}

	if handshake.State != 2 {
		return
	}

	name, uuid, err := minecraftConnection.Login(handshake, server.authenticator, server.protocolVersion, server.onlineMode)

	if err != nil {
		minecraftConnection.WritePacket(&protocol.LoginDisconnect{(&message.Message{Text: err.Error(), Color: message.Red}).JSONString()})
		return
	}

	log.Println("Successfully authenticated " + name + " (" + uuid + ")!")

	player := entity.NewPlayer(name, uuid, minecraftConnection)

	player.Join()
}
