package kafka

import (
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

type SyncProducer interface {
	SendMessage(topic string, content string) error
	SendMessageWithKey(topic string, key string, content string) error
}

type syncProducer struct {
	Producer  sarama.SyncProducer
	closeOnce sync.Once
}

func NewSyncProducer(cfg *ProducerConfig) (SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(strings.Split(cfg.Brokers, ","), cfg.Config)
	if err != nil {
		return nil, err
	}

	p := &syncProducer{
		Producer: producer,
	}

	return p, nil
}

func (p *syncProducer) Close() (err error) {
	p.closeOnce.Do(func() {
		err = p.Producer.Close()
	})
	return
}

func (p *syncProducer) SendMessage(topic string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = int32(-1)
	msg.Value = sarama.ByteEncoder(content)

	_, _, err := p.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *syncProducer) SendMessageWithKey(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = int32(-1)
	msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.ByteEncoder(content)

	_, _, err := p.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
