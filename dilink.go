// dilink 는 웹에서 dilink:// 로 시작하는 프로토콜을 인식하고,
// 관련 응용프로그램으로 URL값을 넘겨주는 프로그램이다.
package main

import (
	"di/dipath"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const RV_WIN = "C:\\Program Files\\Shotgun\\RV-7.0\\bin\\rv.exe"
const RV_lin = "/opt/rv-Linux-x86-64-7.0.0/bin/rv"
const RV_MAC = "/Applications/RV64.app/Contents/MacOS/RV64"
const REGCODE_WIN = `Windows Registry Editor Version 5.00
[HKEY_CLASSES_ROOT\dilink]
@="URL:DIlink Protocol"
"URL Protocol"=""

[HKEY_CLASSES_ROOT\dilink\DefaultIcon]
@="dilink.exe,1"

[HKEY_CLASSES_ROOT\dilink\shell]

[HKEY_CLASSES_ROOT\dilink\shell\open]

[HKEY_CLASSES_ROOT\dilink\shell\open\command]
@="\"\\\\10.0.200.100\\_lustre_INHouse\\Windows\\bin\\dilink.exe\" \"%1\""`

// installDilink는 dilink를 설치하는 함수이다.
func installDilink() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	switch runtime.GOOS {
	case "windows":
		//gen_regcode
		regfile, err := os.Create(user.HomeDir + "\\" + "dilink.reg")
		if err != nil {
			fmt.Fprintf(os.Stderr, "dilink: can't create %s, %s\n", user.HomeDir+"\\"+"dilink.reg", err)
		}
		//윈도우즈라서 캐리지리턴으로 문자열을 바꾸었다.
		if _, err := regfile.Write([]byte(strings.Replace(REGCODE_WIN, "\n", "\r\n", -1))); err != nil {
			fmt.Fprintf(os.Stderr, "dilink: can't save %s, %s\n", user.HomeDir+"\\"+"dilink.reg", err)
		}
		regfile.Close()
		fmt.Println(user.HomeDir + "\\" + "dilink.reg")
		err = exec.Command("cmd", "/C", "start", "", user.HomeDir+"\\"+"dilink.reg").Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Dilink installed for windows.")
		os.Exit(0)
	case "darwin":
		fmt.Println("macOS는 수동으로 브라우저에서 dilink를 설정해야 합니다.")
		os.Exit(1)
	case "linux":
		err := exec.Command("gconftool-2", "--set", "/desktop/gnome/url-handlers/dilink/command", "--type=string", "/lustre/INHouse/CentOS/bin/dilink %s").Run()
		if err != nil {
			log.Fatal(err)
		}
		err = exec.Command("gconftool-2", "--set", "--type=bool", "/desktop/gnome/url-handlers/dilink/enabled", "true").Run()
		if err != nil {
			log.Fatal(err)
		}
		err = exec.Command("gconftool-2", "--set", "--type=bool", "/desktop/gnome/url-handlers/dilink/need-terminal", "false").Run()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Dilink installed for linux.")
		os.Exit(0)
	default:
		fmt.Println("dilink를 설치할 수 없는 OS입니다.")
		os.Exit(1)
	}
}

// Windows 액션
func runWin(scape string) {
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".mov":
		os.Setenv("RV_SUPPORT_PATH", "//10.0.200.100/_lustre_INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		if strings.Contains(scape, ";") {
			var movlist []string
			pathlist := strings.Split(scape, ";")
			for _, i := range pathlist {
				movlist = append(movlist, dipath.Lin2win(i))
			}
			err := exec.Command(RV_WIN, movlist...).Run()
			if err != nil {
				log.Fatal(err)
			}
			return
		}
		err := exec.Command(RV_WIN, dipath.Lin2win(scape)).Run()
		if err != nil {
			log.Fatal(err)
		}
		return
	case ".rv":
		os.Setenv("RV_SUPPORT_PATH", "//10.0.200.100/_lustre_INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		err := exec.Command(RV_WIN, dipath.Lin2win(scape)).Run()
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

// macOS 액션
func runMac(scape string) {
	switch strings.ToLower(filepath.Ext(scape)) {
	case ".nk":
		os.Setenv("NUKE_PATH", "/lustre/INHouse/nuke")
		os.Setenv("NUKE_OFX", "/usr/OFX")
		os.Setenv("PYTHONPATH", "/lustre/INHouse/CentOS/python26/lib:/lustre/INHouse/CentOS/python26/lib/python2.6/site-packages")
		err := exec.Command(RV_MAC, scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".mov", ".jpg":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1") // Prores코덱을 위해서 활성화 한다.
		err := exec.Command(RV_MAC, scape).Run()
		if err != nil {
			log.Fatal(err)
		}
	default:
		err := exec.Command("open", scape).Run()
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
		err := exec.Command(RV_lin, imagelist...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".rv":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1")                                                    // Prores코덱을 위해서 활성화 한다.
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		err := exec.Command(RV_lin, dipath.Win2lin(scape)).Run()
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
		err := exec.Command(RV_lin, movlist...).Run()
		if err != nil {
			log.Fatal(err)
		}
	case ".mp4":
		os.Setenv("RV_ENABLE_MIO_FFMPEG", "1")                                                    // Prores코덱을 위해서 활성화 한다.
		os.Setenv("RV_SUPPORT_PATH", "/lustre/INHouse/rv/supportPath")                            // 회사 RV 파이프라인툴을 로딩하기 위해서 필요하다.
		os.Setenv("PKG_CONFIG_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/pkgconfig")          // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("LD_LIBRARY_PATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib")                    // RV플러그인중 OpenCV를 로딩하기 위해서 필요함.
		os.Setenv("PYTHONPATH", "/lustre/INHouse/Tool/opencv/v3.2.0/lib/python2.7/site-packages") // import cv2, import numpy 를 로딩하기 위해서 필요하다.
		err := exec.Command(RV_lin, dipath.Win2lin(scape)).Run()
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
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stdout, "명령를 실행하기 위한 인수가 충분하지 않습니다.\n")
		fmt.Fprintf(os.Stdout, "dilink 설치를 원하시면 터미널에서 'dilink install'이라고 타이핑 해주세요.\n")
		os.Exit(1)
	}
	// 프로토콜 설치
	if flag.Args()[0] == "install" {
		installDilink()
	}
	// dilink 프로토콜이 올바르게 써져있는지 체크함.
	if !strings.HasPrefix(flag.Args()[0], "dilink://") {
		fmt.Fprintf(os.Stdout, "인수가 dilink://로 시작하지 않습니다. 종료합니다.")
		os.Exit(1)
	}
	uri := strings.TrimPrefix(flag.Args()[0], "dilink://")
	// URI를 통해서 문자를 받기 때문에 %3A -> ":", %2F -> "/" 같은 문자가 섞일 수 있다.
	// 이러한 문자를 QueryUnescape 함수를 통해서 1차 정리한다.
	scape, err := url.QueryUnescape(uri)
	if err != nil {
		log.Fatal(err)
	}
	setProjectnShot(scape) //`$PROJECT, $SEQ, $SHOT 설정`
	switch runtime.GOOS {
	case "darwin":
		runMac(scape)
	case "linux":
		runLin(scape)
	case "windows":
		runWin(scape)
	default:
		fmt.Fprintf(os.Stdout, "지원하지 않는 OS입니다.")
		os.Exit(1)
	}
}
