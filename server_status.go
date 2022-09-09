package main

import "encoding/json"

type ServerStatus struct {
	Version            ServerStatusVersion     `json:"version"`
	Players            ServerStatusPlayers     `json:"players"`
	Description        ServerStatusDescription `json:"description"`
	PreviewsChat       bool                    `json:"previewsChat"`
	EnforcesSecureChat bool                    `json:"enforcesSecureChat"`
}

type ServerStatusVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerStatusPlayers struct {
	Max    int   `json:"max"`
	Online int   `json:"online"`
	Sample []any `json:"sample"`
}

type ServerStatusDescription struct {
	Text string `json:"text"`
}

func (ss *ServerStatus) Encode() (string, error) {
	buff, err := json.Marshal(ss)
	return string(buff), err
}
