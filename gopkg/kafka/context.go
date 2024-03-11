package kafka

import (
	"github.com/Shopify/sarama"
	"time"
)

// 消费消息上下文
// groupId string, topic string, partition int32, offset int64, message []byte, timeStamp time.Time, consumerVal *sarama.ConsumerMessage
type ConsumerMessageContext struct {
	groupId         string
	topic           string
	partition       int32
	offset          int64
	message         []byte
	timeStamp       time.Time
	val             *sarama.ConsumerMessage
	consumerSession *ConsumerSession
}

func (ctx *ConsumerMessageContext) reset() {
	ctx.groupId = ""
	ctx.topic = ""
	ctx.partition = 0
	ctx.offset = 0
	ctx.message = nil
	ctx.timeStamp = time.Time{}
	ctx.val = nil
	ctx.consumerSession = nil
}

func (ctx *ConsumerMessageContext) GetGroupId() string {
	return ctx.groupId
}
func (ctx *ConsumerMessageContext) GetTopic() string {
	return ctx.topic
}
func (ctx *ConsumerMessageContext) GetPartition() int32 {
	return ctx.partition
}
func (ctx *ConsumerMessageContext) GetOffset() int64 {
	return ctx.offset
}
func (ctx *ConsumerMessageContext) GetMessage() []byte {
	return ctx.message
}
func (ctx *ConsumerMessageContext) GetMessageString() string {
	return string(ctx.message)
}

func (ctx *ConsumerMessageContext) GetTimeStamp() time.Time {
	return ctx.timeStamp
}
func (ctx *ConsumerMessageContext) GetVal() *sarama.ConsumerMessage {
	return ctx.val
}
func (ctx *ConsumerMessageContext) GetSession() *ConsumerSession {
	return ctx.consumerSession
}

type ConsumerSession struct {
	session   sarama.ConsumerGroupSession
	message   *sarama.ConsumerMessage
	isAutoAck bool
}

func NewConsumerSession(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage, isAutoAck bool) *ConsumerSession {
	return &ConsumerSession{
		session:   session,
		message:   message,
		isAutoAck: isAutoAck,
	}
}

func (session *ConsumerSession) Ack() {
	if session.session != nil && !session.isAutoAck {
		session.session.MarkMessage(session.message, "")
		session.session.Commit()
	}
}

func (session *ConsumerSession) IsAutoAck() bool {
	return session.isAutoAck
}

func (session *ConsumerSession) GetSession() sarama.ConsumerGroupSession {
	return session.session
}
func (session *ConsumerSession) GetMessage() *sarama.ConsumerMessage {
	return session.message
}
