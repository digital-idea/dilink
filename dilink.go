package main

// dilink 는 웹에서 dilink:// 로 시작하는 URL을 인식하고,
// dilink 명령어에 URL 값을 넘겨 관련 응용프로그램을 실행하는 프로그램이다.

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/digital-idea/dipath"
)

const (
	rvWindowsAppPath = "C:\\Program Files\\Shotgun\\RV-7.0\\bin\\rv.exe"
	rvLinuxAppPath   = "/opt/rv-Linux-x86-64-7.0.0/bin/rv"
	rvMacosAppPath   = "/Applications/RV64.app/Contents/MacOS/RV64"
	protocol         = "dilink://"
)

// Windows 액션
func runWin(scape string) {
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

// Linux 액션
func runLin(scape string) {
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".nk":
		// 회사 셋팅에서 사용자 .bashrc에 보면 IP팀이 umask 0002라고 설정해놓았다.
		// dilink를 통해서 뉴크를 실행하기 때문에 dilink 도 umask 설정이 필요하다.
		// 이렇게 설정이되어야 뉴크실행후 뉴크가 만드는 폴더에 대해서 권한문제가 발생하지 않는다.
		syscall.Umask(0002) // 윈도우는 지원 안함.
		os.Setenv("NUKE_PATH", "/lustre/INHouse/nuke")
		os.Setenv("NUKE_OFX", "/usr/OFX")
		os.Setenv("OPTICAL_FLARES_LICENSE_SERVER_IP", "10.0.99.15")
		os.Setenv("BROWSER", "firefox")
		os.Setenv("NUKE_FONT_PATH", "/lustre2/Digitalidea_source/2d_team_source/font")
		os.Setenv("PYTHONPATH", "/lustre/INHouse/CentOS/python26/lib:/lustre/INHouse/CentOS/python26/lib/python2.6/site-packages")
		os.Setenv("NUKE_USE_FAST_ALLOCATOR", "1")
		err := exec.Command("/usr/local/Nuke10.0v5/Nuke10.0", "--nukex", scape).Run()
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
		err := exec.Command("gnome-terminal", "-x", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".hip":
		syscall.Umask(0002) // 윈도우는 지원 안함.
		err := exec.Command("gnome-terminal", "-x", "h", scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".usd", ".usda":
		syscall.Umask(0002) // 윈도우는 지원 안함.
		err := exec.Command("gnome-terminal", "-x", "uview", scape).Run()
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
		runLin(scape)
	case "windows":
		setProjectnShot(scape) //`digitalidea $PROJECT, $SEQ, $SHOT 설정`
		runWin(scape)
	default:
		fmt.Fprintf(os.Stdout, "지원하지 않는 OS입니다.\n")
		os.Exit(1)
	}
}
