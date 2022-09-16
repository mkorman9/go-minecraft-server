package main

import "encoding/json"

type ChatMessage struct {
	Text            string         `json:"text"`
	IsBold          bool           `json:"bold,omitempty"`
	IsItalic        bool           `json:"italic,omitempty"`
	IsUnderlined    bool           `json:"underlined,omitempty"`
	IsStrikethrough bool           `json:"strikethrough,omitempty"`
	IsObfuscated    bool           `json:"obfuscated,omitempty"`
	FontName        string         `json:"font,omitempty"`
	Extra           []*ChatMessage `json:"extra,omitempty"`
}

const (
	FontDefault = "minecraft:default"
)

func NewChatMessage(text string) *ChatMessage {
	return &ChatMessage{
		Text:     text,
		FontName: FontDefault,
	}
}

func (cm *ChatMessage) Bold() *ChatMessage {
	cm.IsBold = true
	return cm
}

func (cm *ChatMessage) Italic() *ChatMessage {
	cm.IsItalic = true
	return cm
}

func (cm *ChatMessage) Underlined() *ChatMessage {
	cm.IsUnderlined = true
	return cm
}

func (cm *ChatMessage) Strikethrough() *ChatMessage {
	cm.IsStrikethrough = true
	return cm
}

func (cm *ChatMessage) Obfuscated() *ChatMessage {
	cm.IsObfuscated = true
	return cm
}

func (cm *ChatMessage) Font(font string) *ChatMessage {
	cm.FontName = font
	return cm
}

func (cm *ChatMessage) Append(msg *ChatMessage) *ChatMessage {
	cm.Extra = append(cm.Extra, msg)
	return cm
}

func (cm *ChatMessage) Encode() string {
	encoded, _ := json.Marshal(cm)
	return string(encoded)
}
