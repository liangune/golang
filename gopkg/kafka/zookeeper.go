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

func GetKafkaAddress(zkAdress string, chrootPath string) (string, error) {
	zkConn, _, err := zk.Connect(strings.Split(zkAdress, ","), 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("zookeeper connect is failed, %v", err)
	}
	defer zkConn.Close()

	ids, _, err := zkConn.Children(combineKakfaPath(chrootPath, "/brokers/ids"))
	if err != nil {
		return "", fmt.Errorf("zookeeper conn get get brokers ids failed, %v", err)
	}
	kakfaaddress := make([]string, len(ids))
	zkkafka := Kafka{}
	for i, id := range ids {
		result, _, err := zkConn.Get(combineKakfaPath(chrootPath, "/brokers/ids/", id))
		if err != nil {
			return "", fmt.Errorf("zookeeper conn get brokers id[%d] info failed, %v", id, err)
		}
		err = json.Unmarshal(result, &zkkafka)
		if err != nil {
			return "", fmt.Errorf("json unmarshal str[%s] to struct fail, %v", string(result), err)
		}
		kakfaaddress[i] = fmt.Sprintf("%s:%d", zkkafka.Host, zkkafka.Port)
	}

	adress := strings.Join(kakfaaddress, ",")
	return adress, nil
}

func combineKakfaPath(chrootPath string, path ...string) string {
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

func GetKafkaBrokers(zkAdress string, chrootPath string) (map[int]*KakfkaBroker, error) {
	zkConn, _, err := zk.Connect(strings.Split(zkAdress, ","), 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("zookeeper connect is failed, %v", err)
	}
	defer zkConn.Close()

	KafkaBrokers := make(map[int]*KakfkaBroker, 0)
	ids, _, err := zkConn.Children(combineKakfaPath(chrootPath, "/brokers/ids"))
	if err != nil {
		return nil, fmt.Errorf("zookeeper conn get brokers ids failed, %v", err)
	}
	for _, id := range ids {
		broker, _, err := zkConn.Get(combineKakfaPath(chrootPath, "/brokers/ids", id))
		if err != nil {
			return nil, fmt.Errorf("zookeeper conn get brokers id[%d] info failed, %v", id, err)
		}
		v := &KakfkaBroker{}
		err = json.Unmarshal(broker, v)
		if err != nil {
			return nil, fmt.Errorf("json unmarshal str[%s] to struct fail, %v", string(broker), err)
		}
		intid, _ := strconv.Atoi(id)
		v.Id = intid
		KafkaBrokers[intid] = v
	}
	return KafkaBrokers, nil
}
