module github.com/prairir/JobProtocol/Jobs

go 1.15

require (
	github.com/bachittle/ping-go v1.1.10
	github.com/bachittle/ping-go/pinger v0.0.0-20201204235757-2c8f56eb655c
	github.com/bachittle/ping-go/utils v0.0.0-20201204235757-2c8f56eb655c
	github.com/google/gopacket v1.1.19
	github.com/prairir/JobProtocol/Globals v0.0.0-20201129181807-5defa380f9c9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace github.com/prairir/JobProtocol/Globals => ../Globals
