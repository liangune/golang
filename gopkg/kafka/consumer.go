package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"strings"
)

type GroupConsumer struct {
	Consumer        *cluster.Consumer
	Brokers         string
	Topic           string
	GroupID         string
	ClientID        string
	isAutoCommit    bool
	AutoResetOffset string
}

func NewGroupConsumer(cfg *ConsumerConfig) (*GroupConsumer, error) {
	consumer, err := cluster.NewConsumer(strings.Split(cfg.Brokers, ","), cfg.GroupID, []string{cfg.Topic}, cfg.ClusterConfig)
	if err != nil {
		return nil, err
	}

	c := &GroupConsumer{
		Consumer:     consumer,
		Brokers:      cfg.Brokers,
		Topic:        cfg.Topic,
		isAutoCommit: cfg.ClusterConfig.Consumer.Offsets.AutoCommit.Enable,
	}

	if cfg.ClientID != "" {
		cfg.Config.ClientID = cfg.ClientID
	}

	if cfg.AutoResetOffset == AutoResetOffsetEarliest {
		cfg.ClusterConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		cfg.ClusterConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	return c, nil
}

func (c *GroupConsumer) Close() error {
	return c.Consumer.Close()
}

func (c *GroupConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.Consumer.Messages()
}

func (c *GroupConsumer) Notifications() <-chan *cluster.Notification {
	return c.Consumer.Notifications()
}

func (c *GroupConsumer) Errors() <-chan error {
	return c.Consumer.Errors()
}

func (c *GroupConsumer) GetBrokers() string {
	return c.Brokers
}

func (c *GroupConsumer) GetGroupId() string {
	return c.GroupID
}

func (c *GroupConsumer) GetTopic() string {
	return c.Topic
}

func (c *GroupConsumer) IsAutoCommit() bool {
	return c.isAutoCommit
}

// PartitionConsumer can not set GroupID
type PartitionConsumer struct {
	Consumer             sarama.PartitionConsumer
	ParentConsumer       sarama.Consumer
	Brokers              string
	Topic                string
	GroupID              string
	Partition            int32
	ClientID             string
	isAutoCommit         bool
	AutoResetOffset      string
	OffsetBeginTimestamp int64 // given time in seconds
	OffsetEndTimestamp   int64 // given time in seconds
}

func NewPartitionConsumer(cfg *ConsumerConfig) (*PartitionConsumer, error) {
	client, err := sarama.NewClient(strings.Split(cfg.Brokers, ","), cfg.Config)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	c, err := sarama.NewConsumer(strings.Split(cfg.Brokers, ","), cfg.Config)
	if err != nil {
		return nil, err
	}

	var partitionConsumer sarama.PartitionConsumer
	if cfg.AutoResetOffset == AutoResetOffsetTimestamp {
		// given time (in milliseconds) on the topic/partition
		nextOffset, err := client.GetOffset(cfg.Topic, cfg.Partition, cfg.OffsetBeginTimestamp*1000)
		if err != nil {
			return nil, err
		}

		partitionConsumer, err = c.ConsumePartition(cfg.Topic, cfg.Partition, nextOffset)
		if err != nil {
			return nil, err
		}

	} else if cfg.AutoResetOffset == AutoResetOffsetEarliest {
		partitionConsumer, err = c.ConsumePartition(cfg.Topic, cfg.Partition, sarama.OffsetOldest)
		if err != nil {
			return nil, err
		}
	} else {
		partitionConsumer, err = c.ConsumePartition(cfg.Topic, cfg.Partition, sarama.OffsetNewest)
		if err != nil {
			return nil, err
		}
	}

	pcs := PartitionConsumer{
		Consumer:             partitionConsumer,
		ParentConsumer:       c,
		Topic:                cfg.Topic,
		Partition:            cfg.Partition,
		GroupID:              "",
		isAutoCommit:         cfg.Config.Consumer.Offsets.AutoCommit.Enable,
		AutoResetOffset:      cfg.AutoResetOffset,
		OffsetBeginTimestamp: cfg.OffsetBeginTimestamp,
		OffsetEndTimestamp:   cfg.OffsetEndTimestamp,
	}

	if cfg.ClientID != "" {
		cfg.Config.ClientID = cfg.ClientID
	}

	return &pcs, nil
}

func (pcs *PartitionConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return pcs.Consumer.Messages()
}

func (pcs *PartitionConsumer) Errors() <-chan *sarama.ConsumerError {
	return pcs.Consumer.Errors()
}

func (pcs *PartitionConsumer) Close() error {
	pcs.Consumer.AsyncClose()
	err := pcs.ParentConsumer.Close()
	return err
}

func (pcs *PartitionConsumer) GetBrokers() string {
	return pcs.Brokers
}

func (pcs *PartitionConsumer) GetGroupId() string {
	return pcs.GroupID
}

func (pcs *PartitionConsumer) GetTopic() string {
	return pcs.Topic
}

func (pcs *PartitionConsumer) IsAutoCommit() bool {
	return pcs.isAutoCommit
}
