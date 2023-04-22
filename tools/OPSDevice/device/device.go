package device

// 设备状态
const (
	DeviceStatusEable     = 1 // 正常
	DeviceStatusDisable   = 2 // 停用
	DeviceStatusBreakdown = 3 // 故障
	DeviceStatusGiveout   = 4 // 报废
)

const (
	AccessNetworkTypePolice      = 1  // 公安信息网
	AccessNetworkTypeVideo       = 2  // 视频专网
	AccessNetworkTypeTraffic     = 3  // 交管专网
	AccessNetworkTypeInternet    = 4  // 互联网
	AccessNetworkTypeGovernment  = 5  // 政务网
	AccessNetworkTypeEnterprise  = 6  // 企事业单位专网
	AccessNetworkTypePrivate     = 7  // 有局域网
	AccessNetworkTypeRecreation  = 8  // 旅业/娱乐场所专网
	AccessNetworkTypeOtherPolice = 9  // 他公安业务专网
	AccessNetworkTypeOther       = 99 // 其他
)

const (
	IsOnline  = "1" // 在线
	IsOffline = "2" // 离线
	IsOther   = "9" // 其他
)

// 连接协议
const (
	ProtocolNationalStandard = 0  // 28181 国标
	ProtocolSDK              = 1  // SDK
	ProtocolONVIF            = 2  // ONVIF
	Protocol1400             = 3  // GAT 1400标准
	ProtocolOther            = 99 // 其他
)

func GetProtocolName(protocol int) string {
	switch protocol {
	case ProtocolNationalStandard:
		return "28181国标"
	case ProtocolSDK:
		return "SDK"
	case ProtocolONVIF:
		return "ONVIF"
	case Protocol1400:
		return "GAT1400"
	case ProtocolOther:
		return "其他"
	default:
		return "其他"
	}
}

type Device struct {
	DeviceId         string
	DeviceName       string
	IpAddr           string
	Protocol         int
	IsOnline         string
	StatusUpdateTime string
}
