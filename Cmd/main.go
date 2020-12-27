package main

import (
	"bufio"
	"fmt"
	"github.com/prairir/JobProtocol/Creator"
	"github.com/prairir/JobProtocol/Seeker"
	"os"
	"strings"
	"time"
)

func main() {
	var arg string
	if len(os.Args) > 1 {
		arg = os.Args[1]
	}
	switch strings.ToLower(arg) {
	case "creator":
		in := make(chan string)
		out := make(chan map[string]string, 1)
		go creator.RunCreator(in, out, nil)
		//go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Enter a query (ex: JOB EQN 2+2)")
			fmt.Print("> ")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			in <- text
			//}
			//}()
			//for {
			res := <-out
			fmt.Println("result:", res)
		}
	case "seeker":
		seeker.Cmd()
		break
	}
}
