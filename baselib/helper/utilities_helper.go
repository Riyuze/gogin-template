package helper

import (
	"runtime"
	"strconv"

	"github.com/pkg/errors"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CatchErr(err error) error {
	pc, filename, line, _ := runtime.Caller(1)
	return errors.Wrap(err, runtime.FuncForPC(pc).Name()+";"+filename+";"+strconv.Itoa(line))
}
