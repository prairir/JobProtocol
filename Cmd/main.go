package main

import (
	"github.com/prairir/JobProtocol/Creator"
	"github.com/prairir/JobProtocol/Seeker"
	"os"
	"strings"
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch strings.ToLower(arg) {
	case "creator":
		creator.RunCreator(nil, nil, nil, nil)
		break
	case "seeker":
		seeker.Seeker()
		break
	}
}
