package protocol

type Ping struct {
	Version     PingVersion `json:"version"`
	Players     PingPlayers `json:"players"`
	Description string      `json:"description"`
	Favicon     string      `json:"favicon,omitempty"`
}

type PingVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type PingPlayers struct {
	Max    int          `json:"max"`
	Online int          `json:"online"`
	Sample []PingPlayer `json:"sample,omitempty"`
}

type PingPlayer struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
