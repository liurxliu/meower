package event

import (
	"bytes"
	"encoding/gob"

	"github.com/liurxliu/meower/schema"
	"github.com/nats-io/go-nats"
)

type NatsEventStore struct {
	nc                      *nats.Conn
	meowCreatedSubscription *nats.Subscription
	meowCreatedChan         chan MeowCreateMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (e *NatsEventStore) Close() {
	if e.nc != nil {
		e.nc.Close()
	}
	if e.meowCreatedSubscription != nil {
		e.meowCreatedSubscription.Unsubscribe()
	}
	close(e.meowCreatedChan)
}

func (e *NatsEventStore) PublishMeowCreated(meow schema.Meow) error {
	m := MeowCreateMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := e.writeMessage(&m)
	if err != nil {
		return err
	}
	return e.nc.Publish(m.Key(), data)
}

func (mq *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *NatsEventStore) OnMeowCreated(f func(MeowCreateMessage)) (err error) {
	m := MeowCreateMessage{}
	e.meowCreatedSubscription, err = e.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		e.readMessage(msg.Data, &m)
		f(m)
	})
	return
}

func (mq *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

func (e *NatsEventStore) SubscribeMeowCreated() (<-chan MeowCreateMessage, error) {
	m := MeowCreateMessage{}
	e.meowCreatedChan = make(chan MeowCreateMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	e.meowCreatedSubscription, err = e.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case msg := <-ch:
				e.readMessage(msg.Data, &m)
				e.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan MeowCreateMessage)(e.meowCreatedChan), nil
}
