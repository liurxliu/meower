package event

import (
	"github.com/liurxliu/meower/schema"
)

type EventStore interface {
	Close()
	PublishMeowCreated(meow schema.Meow) error
	SubscribeMeowCreated() (<-chan MeowCreateMessage, error)
	OnMeowCreated(f func(MeowCreateMessage)) error
}

var impl EventStore

func SetEventStore(es EventStore) {
	impl = es
}

func Close() {
	impl.Close()
}

func PublishMeowCreated(meow schema.Meow) error {
	return impl.PublishMeowCreated(meow)
}

func SubscribeMeowCreated() (<-chan MeowCreateMessage, error) {
	return impl.SubscribeMeowCreated()
}

func OnMeowCreated(f func(MeowCreateMessage)) error {
	return impl.OnMeowCreated(f)
}
