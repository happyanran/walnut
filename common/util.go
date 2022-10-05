package common

import (
	"fmt"
	"runtime"

	"golang.org/x/crypto/bcrypt"
)

type Utils struct{}

func NewUtil() *Utils {
	return &Utils{}
}

func (u Utils) PwdEnrypt(str string) string {
	pwd, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	return string(pwd)
}

func (u Utils) PwdCheck(pwd string, checkpwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(checkpwd)); err != nil {
		return false
	}

	return true
}

func (u Utils) GetCodeLine() string {
	var str string
	for i := 1; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		pcName := runtime.FuncForPC(pc).Name()
		if ok && pcName != "runtime.main" {
			str += fmt.Sprintln(fmt.Sprintf("%s:%d:%s", file, line, pcName))
		} else {
			return str
		}
	}

	return str
}
