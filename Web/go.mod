module github.com/prairir/JobProtocol/Web

go 1.15

replace github.com/prairir/JobProtocol/Creator => ../Creator

replace github.com/prairir/JobProtocol/Globals => ../Globals

require (
	github.com/prairir/JobProtocol/Creator v0.0.0-00010101000000-000000000000
	github.com/prairir/JobProtocol/Globals v0.0.0-00010101000000-000000000000
)
