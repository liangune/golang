package kafka

import (
	"github.com/Shopify/sarama"
	"strings"
	"time"
)

type SyncProducer interface {
	SendMessage(topic string, content string) error
}

type syncProducer struct {
	Producer sarama.SyncProducer
}

func NewSyncProducer(host string) (SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second

	producer, err := sarama.NewSyncProducer(strings.Split(host, ","), config)
	if err != nil {
		return nil, err
	}

	p := &syncProducer{
		Producer: producer,
	}

	return p, nil
}

func (p *syncProducer) SendMessage(topic string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = int32(-1)
	//msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.ByteEncoder(content)

	_, _, err := p.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
