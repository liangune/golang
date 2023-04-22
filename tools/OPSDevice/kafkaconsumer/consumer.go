package kafkaconsumer

import (
	"encoding/json"
	"go/gopkg/kafka"
	"go/gopkg/logger/vglog"
	"go/tools/OPSDevice/device"
)

type Consumer struct {
	GroupConsumer *kafka.GroupConsumer
	MessageChan   chan []byte
}

var DefaultConsumer *Consumer

const DefaultMessageChanSize = 1000

func NewConsumer(address string, topic string, groupID string) *Consumer {
	groupConsumer, err := kafka.NewGroupConsumer(address, topic, groupID, ConsumerCallback)
	if err != nil {
		vglog.Error("kafka.NewGroupConsumer fail: %v", err)
		return nil
	}
	c := &Consumer{
		GroupConsumer: groupConsumer,
		MessageChan:   make(chan []byte, DefaultMessageChanSize),
	}
	go c.HandleMessage()
	return c
}

func ConsumerCallback(data []byte) {
	if DefaultConsumer != nil {
		DefaultConsumer.AddMessage(data)
	}
}

func (c *Consumer) AddMessage(data []byte) {
	c.MessageChan <- data
}

func (c *Consumer) HandleMessage() {
	for {
		select {
		case msg := <-c.MessageChan:
			var message Message
			err := json.Unmarshal(msg, &message)
			if err != nil {
				vglog.Error("%v", err)
				break
			}
			if message.Protocol == messageProtocol1400 {
				val, ok := device.DefaultDeviceManager.GATMap[message.DeviceID]
				if ok {
					val.IsOnline = message.Type
					val.StatusUpdateTime = message.LocalTime
					device.DefaultDeviceManager.GATMap[message.DeviceID] = val
				}
			} else if message.Protocol == messageProtocolSDK {
				val, ok := device.DefaultDeviceManager.SDKMap[message.DeviceID]
				if ok {
					val.IsOnline = message.Type
					val.StatusUpdateTime = message.LocalTime
					device.DefaultDeviceManager.SDKMap[message.DeviceID] = val
				}
			}
		}
	}
}
