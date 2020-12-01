module github.com/prairir/JobProtocol/Jobs

go 1.15

require (
	github.com/bachittle/ping-go v1.1.4
	github.com/bachittle/ping-go/pinger v0.0.0-20201130210025-5fc73c038b76
	github.com/bachittle/ping-go/utils v0.0.0-20201130210025-5fc73c038b76
	github.com/google/gopacket v1.1.19
	github.com/prairir/JobProtocol/Globals v0.0.0-20201129181807-5defa380f9c9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace github.com/prairir/JobProtocol/Globals => ../Globals
