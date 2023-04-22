package kafkaconsumer

/*
{
  "KeepaliveObject": {
    "DeviceID": "37140000005035300009",
    "LocalTime": "20220617111932",
    "type":"1",
    "protocol": "1400"
  }
}
*/

const (
	messageProtocol1400 = "1400"
	messageProtocolSDK  = "SDK"
)

type Message struct {
	DeviceID  string `json:"DeviceID"`
	LocalTime string `json:"LocalTime"`
	Type      string `json:"type"`
	Protocol  string `json:"protocol"`
}
