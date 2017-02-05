package server

type Settings struct {
	Motd       string `yaml:"motd"`
	MaxPlayers int32  `yaml:"max_players"`
	OnlineMode bool   `yaml:"online_mode"`
}
