package globals

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// functional utils for the JobProtocol

// GetHeader returns the header portion of the JobProtocol query
// ex:
// JOB EQN 2+2
// the header portion is JOB EQN
func GetHeader(result string) (string, error) {
	r, err := regexp.Compile("^[a-zA-Z]+(\\s)*[a-zA-Z]+")
	if err != nil {
		return result, err
	}
	// the header info of the message is here
	// ex: for "JOB EQN 2+2", the queryStr is "JOB EQN"
	rMatchList := r.FindStringSubmatch(result)
	if len(rMatchList) > 0 {
		return strings.ToUpper(rMatchList[0]), nil
	}
	return "", errors.New("invalid query")
}

// FatalErrorCheck exits the program with exit code 1 if there is an error.
func FatalErrorCheck(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
