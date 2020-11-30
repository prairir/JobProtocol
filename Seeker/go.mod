module github.com/prairir/JobProtocol/Seeker

replace github.com/prairir/JobProtocol/Globals => ../Globals

go 1.15

require (
	github.com/Knetic/govaluate v3.0.0+incompatible
	github.com/google/gopacket v1.1.19
	github.com/prairir/JobProtocol/Globals v0.0.0-20201129181807-5defa380f9c9
	github.com/prairir/JobProtocol/Jobs v0.0.0-20201129215838-d0a9081829b5
)

replace github.com/prairir/JobProtocol/Globals => ../Globals