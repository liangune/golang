package kafka

import (
	"fmt"
	"testing"
)

func TestKafka(t *testing.T) {
	zookeeperAddress := "192.168.48.128:2181,192.168.48.128:2182,192.168.48.128:2183"
	//zookeeperAddress := "192.168.33.219:2181"
	kafkaAddress, _ := GetKafkaAddress(zookeeperAddress, "")
	fmt.Println(kafkaAddress)
}
