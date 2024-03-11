package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"strings"
	"sync"
)

type ConsumerGroupHandle func(messages *sarama.ConsumerMessage)

type ConsumerGroupHandler struct {
	handle ConsumerGroupHandle
}

func (h ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.handle(msg)
		session.MarkMessage(msg, "")
	}

	return nil
}

type ConsumerGroup struct {
	consumer        sarama.ConsumerGroup
	Brokers         string
	Topic           string
	GroupID         string
	ClientID        string
	isAutoCommit    bool
	AutoResetOffset string
	listener        ConsumerListener
	errorListener   ConsumerErrorListener
	closer          chan struct{}
	closed          chan struct{}
	closerError     chan struct{}
	closedError     chan struct{}
	closeOnce       sync.Once
	ctx             context.Context
	cancelFunc      context.CancelFunc
}

func NewConsumerGroup(cfg *ConsumerConfig) (*ConsumerGroup, error) {
	consumer, err := sarama.NewConsumerGroup(strings.Split(cfg.Brokers, ","), cfg.GroupID, cfg.Config)
	if err != nil {
		return nil, err
	}

	c := &ConsumerGroup{
		consumer:      consumer,
		Brokers:       cfg.Brokers,
		Topic:         cfg.Topic,
		isAutoCommit:  cfg.Config.Consumer.Offsets.AutoCommit.Enable,
		listener:      cfg.Listener,
		errorListener: cfg.ErrorListener,
		closer:        make(chan struct{}),
		closed:        make(chan struct{}),
		closerError:   make(chan struct{}),
		closedError:   make(chan struct{}),
	}

	if cfg.ClientID != "" {
		cfg.Config.ClientID = cfg.ClientID
	}

	if cfg.AutoResetOffset == AutoResetOffsetEarliest {
		cfg.Config.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg.Config.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	return c, nil
}

func (c *ConsumerGroup) Start() {
	go c.Consume()
	go c.ReturnError()
}

func (c *ConsumerGroup) Consume() {
	defer close(c.closed)
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())

	for {
		select {
		case <-c.closer:
			return
		default:
			topics := []string{c.Topic}
			handler := ConsumerGroupHandler{
				handle: c.RecvMessage,
			}
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			err := c.consumer.Consume(c.ctx, topics, handler)
			if err != nil {
				println(err)
			}
		}
	}
}

func (c *ConsumerGroup) RecvMessage(msg *sarama.ConsumerMessage) {
	ctx := &ConsumerMessageContext{
		groupId:   c.GroupID,
		topic:     c.Topic,
		partition: msg.Partition,
		offset:    msg.Offset,
		message:   msg.Value,
		timeStamp: msg.Timestamp,
		val:       msg,
	}
	c.listener(ctx)
}

func (c *ConsumerGroup) Close() (err error) {
	c.closeOnce.Do(func() {
		close(c.closerError)
		<-c.closedError
		c.cancelFunc()
		close(c.closer)
		<-c.closed
		err = c.consumer.Close()
	})
	return
}

func (c *ConsumerGroup) ReturnError() {
	defer close(c.closedError)
	for {
		select {
		case err := <-c.consumer.Errors():
			if c.errorListener != nil {
				c.errorListener(err)
			}
		case <-c.closerError:
			return
		}
	}
}
