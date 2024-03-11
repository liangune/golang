package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	fpath "path"
	"strconv"
	"strings"
	"time"
)

func GetKafkaAddress(zkAddress string, chrootPath string) (string, error) {
	zkConn, _, err := zk.Connect(strings.Split(zkAddress, ","), 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("zookeeper connect is failed, %v", err)
	}
	defer zkConn.Close()

	ids, _, err := zkConn.Children(combineKafkaPath(chrootPath, "/brokers/ids"))
	if err != nil {
		return "", fmt.Errorf("zookeeper conn get get brokers ids failed, %v", err)
	}
	brokers := make([]string, len(ids))
	broker := Broker{}
	for i, id := range ids {
		result, _, err := zkConn.Get(combineKafkaPath(chrootPath, "/brokers/ids/", id))
		if err != nil {
			return "", fmt.Errorf("zookeeper conn get brokers id[%d] info failed, %v", id, err)
		}
		err = json.Unmarshal(result, &broker)
		if err != nil {
			return "", fmt.Errorf("json unmarshal str[%s] to struct fail, %v", string(result), err)
		}
		brokers[i] = fmt.Sprintf("%s:%d", broker.Host, broker.Port)
	}

	address := strings.Join(brokers, ",")
	return address, nil
}

func combineKafkaPath(chrootPath string, path ...string) string {
	if chrootPath == "" {
		return fpath.Join(path...)
	}

	result := make([]string, 0)
	result = append(result, chrootPath)
	for _, i := range path {
		result = append(result, i)
	}
	return fpath.Join(result...)
}

func GetKafkaBrokers(zkAddress string, chrootPath string) (map[int]*Broker, error) {
	zkConn, _, err := zk.Connect(strings.Split(zkAddress, ","), 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("zookeeper connect is failed, %v", err)
	}
	defer zkConn.Close()

	brokers := make(map[int]*Broker, 0)
	ids, _, err := zkConn.Children(combineKafkaPath(chrootPath, "/brokers/ids"))
	if err != nil {
		return nil, fmt.Errorf("zookeeper conn get brokers ids failed, %v", err)
	}
	for _, id := range ids {
		broker, _, err := zkConn.Get(combineKafkaPath(chrootPath, "/brokers/ids", id))
		if err != nil {
			return nil, fmt.Errorf("zookeeper conn get brokers id[%d] info failed, %v", id, err)
		}
		v := &Broker{}
		err = json.Unmarshal(broker, v)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal str[%s] to struct fail, %v", string(broker), err)
		}
		nId, _ := strconv.Atoi(id)
		v.Id = nId
		brokers[nId] = v
	}
	return brokers, nil
}
