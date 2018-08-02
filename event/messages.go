package event

import (
	"github.com/liurxliu/time"
)

type Message interface {
	Key() string
}

type MeowCreateMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *MeowCreateMessage) Key() string {
	return "meow.created"
}
