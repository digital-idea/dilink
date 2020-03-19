package main

// dilink 는 웹에서 dilink:// 로 시작하는 URL을 인식하고,
// dilink 명령어에 URL 값을 넘겨 관련 응용프로그램을 실행하는 프로그램이다.

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/digital-idea/dipath"
)

// Linux 함수는 URL로 전달받은 문자를 실행하는 함수이다.
func Linux(scape string) {
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".nk":
		// 회사 셋팅에서 사용자 .bashrc에 보면 IP팀이 umask 0002라고 설정해놓았다.
		// dilink를 통해서 뉴크를 실행하기 때문에 dilink 도 umask 설정이 필요하다.
		// 이렇게 설정이되어야 뉴크실행후 뉴크가 만드는 폴더에 대해서 권한문제가 발생하지 않는다.
		syscall.Umask(0002) // 윈도우는 지원 안함.
		err := exec.Command("gnome-terminal", "-x", "nuke", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".jpg", ".png", ".exr", ".tga", ".psd", ".dpx", ".tif":
		imglist := []string{}
		for _, i := range strings.Split(scape, ";") {
			imglist = append(imglist, dipath.Win2lin(i))
		}
		imgext := []string{".jpg", ".png", ".exr", ".tga", ".psd", ".dpx", ".tif"}
		imagelist := []string{}
		for _, img := range imglist {
			for _, ext := range imgext {
				if !strings.Contains(img, ext) {
					continue
				}
				imagelist = append(imagelist, img)
			}
		}
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		err := exec.Command(rvLinuxAppPath, imagelist...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".rv":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1")                                                    // Prores코덱을 위해서 활성화 한다.
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		err := exec.Command(rvLinuxAppPath, dipath.Win2lin(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".mov":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1")                                                    // Prores코덱을 위해서 활성화 한다.
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		playlist := []string{}
		for _, i := range strings.Split(scape, ";") {
			playlist = append(playlist, dipath.Win2lin(i))
		}
		// 플레이 리스트를 받아서 입체 체크를 한다.
		movlist := []string{}
		isStereo := false
		for _, mov := range playlist {
			cmdlist, hasStereo := ToRvStereo(mov)
			if !hasStereo {
				movlist = append(movlist, mov)
				continue
			}
			// RV는 left, right 미디어를 같은 그룹을 묶을 때 "[,]"를 사용한다.
			movlist = append(movlist, "[")
			movlist = append(movlist, cmdlist...)
			movlist = append(movlist, "]")
			isStereo = true
		}
		if isStereo {
			// RV에서 입체 재생을 위해서는 옵션 마지막에 "-stereo scanline" 옵션 필요함.
			movlist = append(movlist, "-stereo")
			movlist = append(movlist, "scanline")
		}
		err := exec.Command(rvLinuxAppPath, movlist...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".mp4":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1")                                                    // Prores코덱을 위해서 활성화 한다.
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		err := exec.Command(rvLinuxAppPath, dipath.Win2lin(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".avi", ".mkv":
		err := exec.Command("/usr/bin/vlc", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".ttf": // 폰트는 폰트브라우저로 연다.
		err := exec.Command("/usr/bin/gnome-font-viewer", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".abc":
		err := exec.Command("/lustre/INHouse/CentOS/bin/abcview", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".pdf":
		err := exec.Command("/usr/bin/evince", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".blend":
		err := exec.Command("/lustre/Applications/Linux/blender/blender-2.75a-linux-glibc211-x86_64/blender", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".obj":
		err := exec.Command("/lustre/Applications/Linux/blender/blender-2.75a-linux-glibc211-x86_64/blender", "--python", "/lustre/INHouse/blender/python/loadobj.py", "--", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".sh":
		err := exec.Command("mate-terminal", "-e", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".hip":
		syscall.Umask(0002) // 윈도우는 지원 안함.
		err := exec.Command("mate-terminal", "-e", "h", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".ma", ".mb":
		os.Setenv("RMSTREE", "/opt/pixar/RenderManForMaya-21.6")
		os.Unsetenv("MYENV")
		syscall.Umask(0002) // 윈도우는 지원 안함.
		// 차후 2018로 변경한다.
		cmd := "/netapp/INHouse/Tool/Ecosystem/bin/ecosystem.py -t maya2017,usd -r maya"
		err := exec.Command("mate-terminal", "-e", cmd, scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".usd", ".usda":
		syscall.Umask(0002) // 윈도우는 지원 안함.
		err := exec.Command("mate-terminal", "-e", "uview", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		browser := "nautilus"
		// 리눅스 release정보를 가지고 온다.
		out, err := exec.Command("cat", "/etc/redhat-release").Output()
		if err != nil {
			log.Fatal(err)
		}
		release := strings.TrimSuffix(string(out), "\n")
		if strings.Contains(release, "CentOS Linux release 7.2.1511 (Core)") {
			// 회사는 CentOS7에서는 caja를 기본 브라우저로 사용한다.
			browser = "caja"
		}
		err = exec.Command(browser, scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
