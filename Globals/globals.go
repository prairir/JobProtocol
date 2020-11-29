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
