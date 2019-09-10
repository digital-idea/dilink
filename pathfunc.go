package main

import (
	"log"
	"os/user"
	"strings"
)

// Home2Abspath 함수는 ~ 문자로 경로가 시작하면 물리적인 경로로 바꾸어준다.
func Home2Abspath(p string) string {
	if !strings.HasPrefix(p, "~") {
		return p
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + strings.TrimPrefix(p, "~")
}
