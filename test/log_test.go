package test

import (
	"fmt"
	"runtime"
	"testing"
)

func TestFileWithLineNum(t *testing.T) {
	var s string
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		pcName := runtime.FuncForPC(pc).Name()
		if ok && pcName != "runtime.main" {
			s += fmt.Sprintln(fmt.Sprintf("%s:%d:%s", file, line, pcName))
		} else {
			break
		}
	}
	t.Log(s)
}
