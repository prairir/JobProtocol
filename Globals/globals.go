package globals

const (
	CONN_ADDR = "localhost"
	CONN_PORT = 7777
	CONN_TYPE = "tcp"
)

func GetJobNames() []string {
	return []string{
		"EQN",
		"HOSTUP",
		"TCPFLOOD",
		"UDPFLOOD"}
}

func GetTCPPorts() []int {
	return []int{
		25,
		80,
		443,
		20,
		21,
		23,
		143,
		3389,
		22,
		53,
		110}
}

func GetUDPPorts() []int {
	return []int{
		53,
		67,
		68,
		69,
		123,
		137,
		138,
		139,
		161,
		162,
		389,
		636}
}
