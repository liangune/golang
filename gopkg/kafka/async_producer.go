package kafka

import (
	"github.com/Shopify/sarama"
	"strings"
	"sync"
	"sync/atomic"
)

type AsyncProducer struct {
	Producer          sarama.AsyncProducer
	successListener   ProducerSuccessListener
	errorListener     ProducerErrorListener
	closer            chan struct{}
	closed            chan struct{}
	isReturnSuccesses atomic.Bool
	closeOnce         sync.Once
}

func NewAsyncProducer(cfg *ProducerConfig) (*AsyncProducer, error) {
	producer, err := sarama.NewAsyncProducer(strings.Split(cfg.Brokers, ","), cfg.Config)
	if err != nil {
		return nil, err
	}

	p := &AsyncProducer{
		Producer:          producer,
		successListener:   cfg.SuccessListener,
		errorListener:     cfg.ErrorListener,
		closer:            make(chan struct{}),
		closed:            make(chan struct{}),
		isReturnSuccesses: atomic.Bool{},
	}
	if cfg.Config.Producer.Return.Errors || cfg.Config.Producer.Return.Successes {
		p.isReturnSuccesses.Store(true)
		go p.ReturnSuccessOrErrorMessage()
	}

	return p, nil
}

func (p *AsyncProducer) Send() chan<- *sarama.ProducerMessage {
	return p.Producer.Input()
}

func (p *AsyncProducer) Successes() <-chan *sarama.ProducerMessage {
	return p.Producer.Successes()
}

func (p *AsyncProducer) Errors() <-chan *sarama.ProducerError {
	return p.Producer.Errors()
}

func (p *AsyncProducer) Close() (err error) {
	p.closeOnce.Do(func() {
		close(p.closer)
		if p.isReturnSuccesses.Load() {
			<-p.closed
		}

		err = p.Producer.Close()
	})
	return
}

func (p *AsyncProducer) SendMessage(topic string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = int32(-1)
	msg.Value = sarama.ByteEncoder(content)

	p.Producer.Input() <- msg

	return nil
}

func (p *AsyncProducer) SendMessageWithKey(topic string, key string, content string) error {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Partition = int32(-1)
	msg.Key = sarama.StringEncoder(key)
	msg.Value = sarama.ByteEncoder(content)

	p.Producer.Input() <- msg

	return nil
}

func (p *AsyncProducer) ReturnSuccessOrErrorMessage() {
	defer close(p.closed)
	for {
		select {
		case msg := <-p.Producer.Successes():
			if p.successListener != nil {
				p.successListener(msg)
			}
		case msg := <-p.Producer.Errors():
			if p.errorListener != nil {
				p.errorListener(msg)
			}
		case <-p.closer:
			return
		}
	}
}
