package main

import (
	"log"
	"os/user"
	"strings"
)

func Home2Abspath(p string) string {
	if !strings.HasPrefix(p, "~") {
		return p
	}
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(p, "~", usr.HomeDir, 1)
}
