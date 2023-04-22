package ping

type PingCmdSort []PingCmd

func (s PingCmdSort) Len() int {
	return len(s)
}

func (s PingCmdSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s PingCmdSort) Less(i, j int) bool {
	return s[i].PacketLoss > s[j].PacketLoss
}
