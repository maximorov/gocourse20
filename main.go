package main

import (
	"context"
	"fmt"
	"solid/log"
	"solid/notification"
	"solid/sender"
	"solid/sender/transport"
)

type A struct {
	name string
}

func (p *A) Name() string {
	return p.name
}

type B struct {
	A
}

func (p *B) Description() string {
	return p.Name() + `something else`
}

type C struct {
	B
}

func (p *C) Description() string {
	return p.B.Description() + `something else`
}

func main() {
	s := sender.NewSender(
		transport.NewKafkaCreator(),
		transport.NewHttpCreator(),
		transport.NewRabbitMQCreator(),
	)
	_ = s.Send(context.Background(), nil, nil)

	zc := log.ZapConfiger{}
	c := zc.Build(
		log.WithEnableCaller(true),
		log.WithServiceVersion(`v1.0.0`),
		log.WithSomethingElse(),
	)

	ms := &myStruct{}
	notify(ms)
}

func notify(n notification.NotificationTargets) error {
	al := n.AllTargets()
	return fmt.Errorf(`some %v`, al)
}

type myStruct struct {
}

func (m myStruct) AllTargets() []notification.Target {
	return nil
}
