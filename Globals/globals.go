package globals

const (
	// ConnAddr is a string constant specifying where the job creator connects to.
	ConnAddr = "localhost"
	// ConnPort is a string constant specifying the port the job creator runs on
	ConnPort = 7777
	// ConnType specifies the connection type for the job creator (usually just TCP)
	ConnType = "tcp"
)

// GetJobNames is an array of job names for the C&C to use
func GetJobNames() []string {
	return []string{
		"EQN",
		"HOSTUP",
		"TCPFLOOD",
		"UDPFLOOD"}
}

// GetTCPPorts is an array of TCP ports
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
		110,
	}
}
