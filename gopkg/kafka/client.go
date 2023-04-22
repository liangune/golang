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
func GetConsumerGroups(kafkaAdress string) (map[string]string, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V1_0_0_0
	admin, err := sarama.NewClusterAdmin(strings.Split(kafkaAdress, ","), cfg)
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

type KafkaClient struct {
	client sarama.Client
}

func NewKafkaClient(kafkaAdress string) (*KafkaClient, error) {
	cfg := sarama.NewConfig()
	client, err := sarama.NewClient(strings.Split(kafkaAdress, ","), cfg)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{client: client}, nil
}

func (c *KafkaClient) Close() {
	c.client.Close()
}

func (c *KafkaClient) GetTopicOffset(topic string) (*TopicOffset, error) {
	partitions, _ := c.client.Partitions(topic)
	ktopic := &TopicOffset{}
	ktopic.Topic = topic
	for _, i := range partitions {
		offset, err := c.client.GetOffset(topic, i, sarama.OffsetNewest)
		if err != nil {
			// vglog.Error("kakfa client get newest offet error: %v", err)
			continue
		}
		ktopic.Partitions = append(ktopic.Partitions, &PartitionOffset{PartitionID: int32(i), Offset: offset})
		ktopic.Offset += offset
	}

	return ktopic, nil
}

func (c *KafkaClient) GetTopicGroupOffset(group string, topic string) (*TopicGroupOffset, error) {
	OffetMgr, err := sarama.NewOffsetManagerFromClient(group, c.client)
	if err != nil {
		return nil, err
	}
	defer OffetMgr.Close()

	partitions, err := c.client.Partitions(topic)
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
		newestOff, err := c.client.GetOffset(topic, partID, sarama.OffsetNewest)
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

func (c *KafkaClient) GetTopic() ([]string, error) {
	return c.client.Topics()
}
