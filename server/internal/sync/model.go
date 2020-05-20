package sync

import "github.com/normegil/toth/server/internal"

type Model struct {
	Version int  `yaml:"version"`
	Data    Data `yaml:"data"`
}

type Data struct {
	Users []internal.User `yaml:"users,omitempty"`
}
