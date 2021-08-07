package common

import (
	"fmt"
)

func True(cond bool, a ...interface{}) {
	Truef(cond, fmt.Sprint(a...))
}

func False(cond bool, a ...interface{}) {
	Truef(!cond, fmt.Sprint(a...))
}

func Truef(cond bool, format string, a ...interface{}) {
	if !cond {
		if a == nil || len(a) == 0 {
			panic(format)
		} else {
			panic(fmt.Sprintf(format, a...))
		}
	}
}

func Falsef(cond bool, format string, a ...interface{}) {
	Truef(!cond, format, a...)
}

func AssertHere() {
	Truef(false, "CANNOT run here")
}
