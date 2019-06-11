package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MacOS 함수는 URL로 전달받은 문자를 실행하는 함수이다.
func MacOS(scape string) {
	scape = Home2Abspath(scape)
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".nk", ".nknc":
		os.Setenv("NUKE_PATH", Home2Abspath("~/nuke"))
		os.Setenv("NUKE_FONT_PATH", Home2Abspath("~/nuke/font"))
		os.Setenv("OCIO", Home2Abspath("~/OpenColorIO-Configs/aces_1.0.3/config.ocio"))
		// 맥은 인터넷 연결되 되어있을 가능성이 높다. 항상 논커머셜로 실행한다.
		err := exec.Command("/Applications/Nuke11.3v2/Nuke11.3v2.app/Contents/MacOS/Nuke11.3v2", "--nukex", "--nc", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".mov", ".jpg":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1") // Prores 코덱을 위해서 활성화 한다.
		err := exec.Command(rvMacosAppPath, scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".blend":
		err := exec.Command("/Applications/Blender/blender.app/Contents/MacOS/blender", "--python", Home2Abspath("~/blender/init.py"), scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".kra":
		err := exec.Command("/Applications/krita.app/Contents/MacOS/krita", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".xcf":
		err := exec.Command("/Applications/GIMP-2.10.app/Contents/MacOS/gimp", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".ntp":
		os.Setenv("NATRON_PLUGIN_PATH", Home2Abspath("~/natron"))
		err := exec.Command("/Applications/Natron.app/Contents/MacOS/Natron", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".svg":
		err := exec.Command("/Applications/Inkscape.app/Contents/MacOS/Inkscape", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		// 일반적으로 abc, hwp, 키노트등의 포멧은 open 명령어로 잘 작동된다.
		err := exec.Command("open", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
