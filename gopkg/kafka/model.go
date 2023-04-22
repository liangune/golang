package kafka

type Kafka struct {
	Port int
	Host string
}

// {"listener_security_protocol_map":{"PLAINTEXT":"PLAINTEXT"},"endpoints":["PLAINTEXT://192.168.48.128:9092"],"jmx_port":-1,"host":"192.168.48.128","timestamp":"1595312023650","port":9092,"version":4}
type KakfkaBroker struct {
	Id   int `json:",omitempty"`
	Host string
	Port int
}

type Topic struct {
	Name  string
	Group []string
}

type TopicGroupOffset struct {
	Topic      string
	Offset     int64
	LogSize    int64
	Lag        int64
	Group      string
	Partitions []*PartitionGroupOffset
}

type PartitionGroupOffset struct {
	PartitionID int32
	Offset      int64
	LogSize     int64
	Lag         int64
}

type TopicOffset struct {
	Topic      string
	Offset     int64
	Partitions []*PartitionOffset
}

type PartitionOffset struct {
	PartitionID int32
	Offset      int64
}
