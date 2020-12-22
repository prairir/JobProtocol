module github.com/prairir/JobProtocol/Jobs

go 1.15

require (
	github.com/bachittle/ping-go-v2 v1.0.1
	github.com/google/gopacket v1.1.19
	github.com/prairir/JobProtocol/Globals v0.0.0-20201129181807-5defa380f9c9
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace github.com/prairir/JobProtocol/Globals => ../Globals
