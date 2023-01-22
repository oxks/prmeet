package er

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func ErrorPrint(err error) {
	if err != nil {
		log.Println(err)
	}
}

// u.Validate() returns lowercase properties, so below the workaround
func FixValidationResult(err error) error {
	// string error
	serr := fmt.Sprintf("%v", err)
	//split string into sentences
	sArr := strings.SplitAfter(serr, "; ")

	for key, value := range sArr {
		f := value[0:1]
		f = strings.ToUpper(f)
		sArr[key] = f + value[1:]
	}
	serr = strings.Join(sArr, "")
	return errors.New(serr)
}
