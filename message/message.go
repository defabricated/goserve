package message

import (
	"bytes"
	"encoding/json"
)

type Message struct {
	Text          string      `json:"text,omitempty"`
	Translate     string      `json:"translate,omitempty"`
	With          []*Message  `json:"with,omitempty"`
	Extra         []*Message  `json:"extra,omitempty"`
	Bold          *bool       `json:"bold,omitempty"`
	Italic        *bool       `json:"italic,omitempty"`
	Underlined    *bool       `json:"underlined,omitempty"`
	Strikethrough *bool       `json:"strikethrough,omitempty"`
	Obfuscated    *bool       `json:"obfuscated,omitempty"`
	Color         Color       `json:"color,omitempty"`
	ClickEvent    *ClickEvent `json:"clickEvent,omitempty"`
	HoverEvent    *HoverEvent `json:"hoverEvent,omitempty"`
}

var (
	_false = false
	False  = &_false
	_true  = true
	True   = &_true
)

func (m *Message) JSONString() string {
	res, _ := json.Marshal(m)
	return string(res)
}

func (m *Message) String() string {
	var buf bytes.Buffer
	m.string(&buf)
	return buf.String()
}

func (m *Message) string(buf *bytes.Buffer) {
	if m.Text != "" {
		buf.WriteString(m.Text)
	} else {
		//TODO

		panic("Translatable strings cannot be stringified yet")
	}
	if m.Extra != nil {
		for _, e := range m.Extra {
			e.string(buf)
		}
	}
}

type ClickEvent struct {
	Action ClickAction `json:"action"`
	Value  string      `json:"value"`
}

type HoverEvent struct {
	Action HoverAction `json:"action"`
	Value  *Message    `json:"value"`
}

type ClickAction string

const (
	OpenUrl        ClickAction = "open_url"
	OpenFile       ClickAction = "open_file"
	RunCommand     ClickAction = "run_command"
	SuggestCommand ClickAction = "suggest_command"
)

type HoverAction string

const (
	ShowText        HoverAction = "show_text"
	ShowAchievement HoverAction = "show_achievement"
	ShowItem        HoverAction = "show_item"
)

type Color string

const (
	Black       Color = "black"
	DarkBlue    Color = "dark_blue"
	DarkGreen   Color = "dark_green"
	DarkAqua    Color = "dark_aqua"
	DarkRed     Color = "dark_red"
	DarkPurple  Color = "dark_purple"
	Gold        Color = "gold"
	Gray        Color = "gray"
	DarkGray    Color = "dark_gray"
	Blue        Color = "blue"
	Green       Color = "green"
	Aqua        Color = "aqua"
	Red         Color = "red"
	LightPurple Color = "light_purple"
	Yellow      Color = "yellow"
	White       Color = "white"
)
