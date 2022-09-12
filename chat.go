package main

import "encoding/json"

type ChatMessage struct {
	Text          string        `json:"text"`
	Bold          bool          `json:"bold,omitempty"`
	Italic        bool          `json:"italic,omitempty"`
	Underlined    bool          `json:"underlined,omitempty"`
	Strikethrough bool          `json:"strikethrough,omitempty"`
	Obfuscated    bool          `json:"obfuscated,omitempty"`
	Font          string        `json:"font,omitempty"`
	Extra         []ChatMessage `json:"extra,omitempty"`
}

const (
	FontDefault = "minecraft:default"
)

func NewChatMessage(text string) *ChatMessage {
	return &ChatMessage{
		Text: text,
		Font: FontDefault,
	}
}

func (cm *ChatMessage) Encode() string {
	encoded, _ := json.Marshal(cm)
	return string(encoded)
}
