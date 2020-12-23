module github.com/prairir/JobProtocol/Cmd

require (
	github.com/prairir/JobProtocol/Creator v0.0.0-20201223035010-c67d396ff354
	github.com/prairir/JobProtocol/Globals v0.0.0-20201223035010-c67d396ff354
	github.com/prairir/JobProtocol/Seeker v0.0.0-20201223035010-c67d396ff354
)

replace github.com/prairir/JobProtocol/Globals => ../Globals

replace github.com/prairir/JobProtocol/Seeker => ../Seeker

replace github.com/prairir/JobProtocol/Creator => ../Creator

replace github.com/prairir/JobProtocol/Jobs => ../Jobs

go 1.15
