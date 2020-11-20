package globals

const (
	CONN_ADDR = "localhost"
	CONN_PORT = 7777
	CONN_TYPE = "tcp"
)

func GetJobNames() []string {
	return []string {
		"EQN",
		"HOSTUP",
		"TCPFLOOD",
		"UDPFLOOD" }
}