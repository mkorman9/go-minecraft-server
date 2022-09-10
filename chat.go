package main

import "encoding/json"

type ChatMessage struct {
	Text          string        `json:"text"`
	Bold          bool          `json:"bold"`
	Italic        bool          `json:"italic"`
	Underlined    bool          `json:"underlined"`
	Strikethrough bool          `json:"strikethrough"`
	Obfuscated    bool          `json:"obfuscated"`
	Font          string        `json:"font"`
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
