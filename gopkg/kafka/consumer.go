package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"strings"
	"sync"
	"time"
)

// 消费者回调函数
type ConsumerCallback func(data []byte)

func PartitionConsumer(address string, topic string, wg *sync.WaitGroup, callback ConsumerCallback) error {
	consumer, err := sarama.NewConsumer(strings.Split(address, ","), nil)
	if err != nil {
		return err
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			continue
		}
		defer pc.AsyncClose()

		wg.Add(1)

		go func(pc sarama.PartitionConsumer) {
			defer wg.Done()
			for msg := range pc.Messages() {
				callback(msg.Value)
			}
		}(pc)
	}

	wg.Wait()

	consumer.Close()

	return nil
}

type GroupConsumer struct {
	Consumer *cluster.Consumer
	Callback ConsumerCallback
	closer   chan struct{}
	closed   chan struct{}
}

func NewGroupConsumer(address string, topic string, groupID string, callback ConsumerCallback) (*GroupConsumer, error) {
	config := cluster.NewConfig()
	config.Consumer.Offsets.AutoCommit.Interval = 500 * time.Millisecond
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.CommitInterval = 500 * time.Millisecond
	consumer, err := cluster.NewConsumer(strings.Split(address, ","), groupID, []string{topic}, config)
	if err != nil {
		return nil, err
	}

	c := &GroupConsumer{
		Consumer: consumer,
		Callback: callback,
		closer:   make(chan struct{}),
		closed:   make(chan struct{}),
	}

	go c.RecvMesssage()
	return c, nil
}

func (c *GroupConsumer) Close() error {
	close(c.closer)
	<-c.closed
	return c.Consumer.Close()
}

func (c *GroupConsumer) RecvMesssage() {
	defer close(c.closed)
	for {
		select {
		case msg := <-c.Consumer.Messages():
			c.Callback(msg.Value)
		case <-c.closer:
			return
		}
	}
}
