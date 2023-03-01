package transport

import (
	"context"
	"encoding/json"
	"solid/l/sender"
	"solid/notification"
	"solid/o"
	sender2 "solid/sender"
	"strings"
)

type KafkaSettings struct {
	BrokerList string  `json:"broker_list"`
	TopicName  string  `json:"topic_name"`
	User       *string `json:"user"`
	Pass       *string `json:"pass"`
}

func NewKafkaCreator() sender2.Transport {
	return &kafkaCreator{}
}

type kafkaCreator struct {
}

func (s *kafkaCreator) Name() string {
	return `kafka`
}

func (s kafkaCreator) Create(n ol.Notification) (sender.Target, error) {
	settings := KafkaSettings{}
	if err := json.Unmarshal([]byte(n.TargetSettings()), &settings); err != nil {
		return nil, err
	}

	topicName := coreBus.TopicName()
	broker := kafka.NewBroker(n.ID().String(), strings.Split(settings.BrokerList, ","))
	topic, err := kafka.NewTopic(broker, coreBus.TopicConfig{
		TopicName: topicName,
	})
	if err != nil {
		return nil, err
	}

	err = broker.AddTopic(topicName, topic)
	if err != nil {
		return nil, err
	}

	return &kafkaTransport{broker, topicName}, nil
}

type kafkaTransport struct {
	broker coreBus.Broker
	topic  coreBus.TopicName
}

func (s *kafkaTransport) Send(ctx context.Context, n ol.Notification, i notification.ItemToProcess) error {
	t, err := s.broker.Topic(s.topic)
	if err != nil {
		return err
	}

	return t.SendMessage(ctx, o.FormMessage(n, i))
}

func (s *kafkaTransport) Die() error {
	return nil
}
