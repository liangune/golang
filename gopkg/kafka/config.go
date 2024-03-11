package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"time"
)

const (
	AutoResetOffsetNone      = "none"
	AutoResetOffsetEarliest  = "earliest"
	AutoResetOffsetLatest    = "latest"
	AutoResetOffsetTimestamp = "timestamp"
)

type ConsumerConfig struct {
	Brokers              string
	Topic                string
	GroupID              string
	ClientID             string
	Partition            int32
	Listener             ConsumerListener
	ErrorListener        ConsumerErrorListener
	Config               *sarama.Config
	ClusterConfig        *cluster.Config
	AutoResetOffset      string
	OffsetBeginTimestamp int64
	OffsetEndTimestamp   int64
}

func NewConsumerConfig() *ConsumerConfig {
	c := &ConsumerConfig{}
	c.ClusterConfig = cluster.NewConfig()
	c.ClusterConfig.Consumer.Offsets.CommitInterval = 500 * time.Millisecond
	c.ClusterConfig.Consumer.Offsets.AutoCommit.Interval = 500 * time.Millisecond
	c.ClusterConfig.Consumer.Offsets.AutoCommit.Enable = true
	c.ClusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.ClusterConfig.Consumer.Return.Errors = true

	c.Config = sarama.NewConfig()
	c.Config.Consumer.Offsets.AutoCommit.Interval = 500 * time.Millisecond
	c.Config.Consumer.Offsets.AutoCommit.Enable = true
	c.Config.Consumer.Offsets.Initial = sarama.OffsetNewest
	c.Config.Consumer.Return.Errors = true

	c.AutoResetOffset = AutoResetOffsetLatest
	return c
}

type ProducerConfig struct {
	Brokers         string
	SuccessListener ProducerSuccessListener
	ErrorListener   ProducerErrorListener
	Config          *sarama.Config
}

func NewProducerConfig() *ProducerConfig {
	c := &ProducerConfig{}
	c.Config = sarama.NewConfig()
	c.Config.Producer.Return.Successes = true
	c.Config.Producer.Timeout = 5 * time.Second
	c.Config.Producer.Return.Errors = true

	return c
}
