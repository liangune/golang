package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"strings"
)

/*
@brief kafka版本 1.0.0以上
@return map的key是groupID，value为consumer
*/
func GetConsumerGroups(brokers string) (map[string]string, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	admin, err := sarama.NewClusterAdmin(strings.Split(brokers, ","), cfg)
	if err != nil {
		return nil, fmt.Errorf("NewClusterAdmin error message: %s", err.Error())
	}
	defer admin.Close()

	groups, err := admin.ListConsumerGroups()
	if err != nil {
		return nil, fmt.Errorf("ListConsumerGroups error message: %s", err.Error())
	}

	return groups, nil
}

type Client struct {
	saramaClient sarama.Client
}

func NewKafkaClient(brokers string) (*Client, error) {
	cfg := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(brokers, ","), cfg)
	if err != nil {
		return nil, err
	}

	return &Client{saramaClient: client}, nil
}

func (c *Client) Close() {
	c.saramaClient.Close()
}

func (c *Client) GetTopicOffset(topic string, time int64) (*TopicOffset, error) {
	partitions, _ := c.saramaClient.Partitions(topic)
	topicOffset := &TopicOffset{}
	topicOffset.Topic = topic
	for _, i := range partitions {
		offset, err := c.saramaClient.GetOffset(topic, i, sarama.OffsetNewest)
		if err != nil {
			continue
		}
		topicOffset.Partitions = append(topicOffset.Partitions, &PartitionOffset{PartitionID: int32(i), Offset: offset})
		topicOffset.Offset += offset
	}

	return topicOffset, nil
}

func (c *Client) GetTopicGroupOffset(group string, topic string) (*TopicGroupOffset, error) {
	OffetMgr, err := sarama.NewOffsetManagerFromClient(group, c.saramaClient)
	if err != nil {
		return nil, err
	}
	defer OffetMgr.Close()

	partitions, err := c.saramaClient.Partitions(topic)
	if err != nil {
		return nil, err
	}

	logSize := int64(0)
	lag := int64(0)
	offset := int64(0)
	parts := make([]*PartitionGroupOffset, 0)
	for _, partID := range partitions {
		partOffMgr, err := OffetMgr.ManagePartition(topic, partID)
		if err != nil {
			continue
		}
		defer partOffMgr.Close()

		// 下一个被消费的offset, 对于一个topic的每个分区的next offset，消费组必须要处理过消息才有缓存的，不然next offset为初始化的OffsetNewest或者OffsetOldest
		partOff, _ := partOffMgr.NextOffset()

		// 最新的offset
		newestOff, err := c.saramaClient.GetOffset(topic, partID, sarama.OffsetNewest)
		if err != nil {
			continue
		}

		diff := int64(0)
		// 初始化设置为OffsetNewest，则未消费过的分区的消息记录丢弃，消费过消息后则从next offset开始
		if partOff == sarama.OffsetNewest || partOff == sarama.OffsetOldest {
			offset += newestOff
		} else {
			diff = newestOff - partOff
			lag += diff
			offset += partOff
		}
		logSize += newestOff
		p := &PartitionGroupOffset{
			PartitionID: partID,
			LogSize:     newestOff,
			Lag:         diff,
			Offset:      partOff,
		}
		parts = append(parts, p)
	}

	topicPrometheus := &TopicGroupOffset{
		Group:      group,
		Topic:      topic,
		Offset:     offset,
		LogSize:    logSize,
		Lag:        lag,
		Partitions: parts,
	}

	return topicPrometheus, nil
}

func (c *Client) GetTopics() ([]string, error) {
	return c.saramaClient.Topics()
}

func (c *Client) GetPartitions(topic string) ([]int32, error) {
	return c.saramaClient.Partitions(topic)
}
