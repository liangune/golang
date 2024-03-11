package kafka

import "github.com/Shopify/sarama"

// 消费者返回数据监听
type ConsumerListener func(context *ConsumerMessageContext)
type ConsumerErrorListener func(err error)

// 异步生产者回调监听
type ProducerSuccessListener func(message *sarama.ProducerMessage)
type ProducerErrorListener func(message *sarama.ProducerError)
