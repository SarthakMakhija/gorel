package gorel

import "fmt"

func Assert(condition bool, msg string, values ...interface{}) {
	if !condition {
		panic(fmt.Sprintf(msg, values...))
	}
}
