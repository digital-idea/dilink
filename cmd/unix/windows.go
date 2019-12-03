package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/digital-idea/dipath"
)

// Windows 함수는 URL로 전달받은 문자를 실행하는 함수이다.
func Windows(scape string) {
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".mov":
		os.Setenv("RV_SUPPORT_PATH", "//10.0.200.100/_lustre_INHouse/rv/supportPath") // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		if strings.Contains(scape, ";") {
			var movlist []string
			pathlist := strings.Split(scape, ";")
			for _, i := range pathlist {
				movlist = append(movlist, dipath.Lin2win(i))
			}
			err := exec.Command(rvWindowsAppPath, movlist...).Run()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		err := exec.Command(rvWindowsAppPath, dipath.Lin2win(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	case ".rv":
		os.Setenv("RV_SUPPORT_PATH", "//10.0.200.100/_lustre_INHouse/rv/supportPath") // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		err := exec.Command(rvWindowsAppPath, dipath.Lin2win(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		err := exec.Command("cmd", "/C", "start", "", dipath.Lin2win(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
