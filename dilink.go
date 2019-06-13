package main

// dilink 는 웹에서 dilink:// 로 시작하는 URL을 인식하고,
// dilink 명령어에 URL 값을 넘겨 관련 응용프로그램을 실행하는 프로그램이다.

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/digital-idea/dipath"
)

const (
	rvWindowsAppPath = "C:\\Program Files\\Shotgun\\RV-7.0\\bin\\rv.exe"
	rvLinuxAppPath   = "/opt/rv-Linux-x86-64-7.0.0/bin/rv"
	rvMacosAppPath   = "/Applications/RV64.app/Contents/MacOS/RV64"
	protocol         = "dilink://"
)

// 해당 경로를 체크하여 PROJECT, SEQ, SHOT의 환경변수 설정
func setProjectnShot(scape string) {
	var project, seq, shot string
	path := scape
	project, err := dipath.Project(path)
	if err != nil {
		log.Println(err)
	}
	seq, err = dipath.Seq(path)
	if err != nil {
		log.Println(err)
	}
	shot, err = dipath.Shot(path)
	if err != nil {
		log.Println(err)
	}
	if project != "" && seq != "" && shot != "" {
		err = os.Setenv("PROJECT", project)
		if err != nil {
			log.Println(err)
		}
		err = os.Setenv("SEQ", seq)
		if err != nil {
			log.Println(err)
		}
		err = os.Setenv("SHOT", shot)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("dilink: ")
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stdout, "명령를 실행하기 위한 인수가 충분하지 않습니다.\n")
		os.Exit(1)
	}
	// dilink 프로토콜이 올바르게 써져있는지 체크함.
	if !strings.HasPrefix(flag.Args()[0], protocol) {
		fmt.Fprintf(os.Stdout, "인수가 %s 로 시작하지 않습니다. 종료합니다.\n", protocol)
		os.Exit(1)
	}
	uri := strings.TrimPrefix(flag.Args()[0], protocol)
	// URI를 통해서 문자를 받기 때문에 %3A -> ":", %2F -> "/" 같은 문자가 섞일 수 있다.
	// 이러한 문자를 QueryUnescape 함수를 통해서 1차 정리한다.
	scape, err := url.QueryUnescape(uri)
	if err != nil {
		log.Fatal(err)
	}

	switch runtime.GOOS {
	case "darwin":
		MacOS(scape)
	case "linux":
		setProjectnShot(scape) //`digitalidea $PROJECT, $SEQ, $SHOT 설정`
		Linux(scape)
	case "windows":
		setProjectnShot(scape) //`digitalidea $PROJECT, $SEQ, $SHOT 설정`
		Windows(scape)
	default:
		fmt.Fprintf(os.Stdout, "지원하지 않는 OS입니다.\n")
		os.Exit(1)
	}
}
