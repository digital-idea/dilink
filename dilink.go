package main

import (
	"os"
	"fmt"
	"os/exec"
	"os/user"
	"log"
	"strings"
	"flag"
	"net/url"
	"runtime"
	"di/dipath"
)

const RV_win = "C:\\Program Files (x86)\\Tweak\\RV-3.12.20-32\\bin\\rv.exe"
const RV_lin = "/opt/rv-Linux-x86-64-4.0.10/bin/rv"
const RV_osx = "/Applications/RV64.app/Contents/MacOS/RV64"
const REGCODE = `Windows Registry Editor Version 5.00
[HKEY_CLASSES_ROOT\dilink]
@="URL:DIlink Protocol"
"URL Protocol"=""

[HKEY_CLASSES_ROOT\dilink\DefaultIcon]
@="dilink.exe,1"

[HKEY_CLASSES_ROOT\dilink\shell]

[HKEY_CLASSES_ROOT\dilink\shell\open]

[HKEY_CLASSES_ROOT\dilink\shell\open\command]
@="\"\\\\10.0.200.100\\_lustre_INHouse\\Windows\\bin\\dilink.exe\" \"%1\""`

func gethomedir() string {
	user, err := user.Current()
	if err != nil {
		return ""
	}
	return fmt.Sprintf(user.HomeDir)
}

func main() {
	var argstr string
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Fprintf(os.Stdout, "not enough argument.\n")
		fmt.Fprintf(os.Stdout, "if you want install dilink protocol then need typing 'dilink install'\n")
		os.Exit(1)
	}
	//install check
	if flag.Args()[0] == "install" {
		switch runtime.GOOS {
			case "windows": {
				//gen_regcode
				regfile, err := os.Create(gethomedir() +"\\"+ "dilink.reg")
				if err != nil {
					fmt.Fprintf(os.Stderr, "dilink: can't create %s, %s\n", gethomedir() + "\\" + "dilink.reg", err)
				}
				//윈도우즈라서 캐리지리턴으로 문자열을 바꾸었다.
				if _, err := regfile.Write([]byte(strings.Replace(REGCODE, "\n", "\r\n", -1))); err != nil {
					fmt.Fprintf(os.Stderr, "dilink: can't save %s, %s\n", gethomedir() + "\\" + "dilink.reg", err)
				}
				regfile.Close()
				//run_regcode()
				fmt.Println(gethomedir()+"\\"+"dilink.reg")
				exec.Command("cmd", "/C", "start", "", gethomedir()+"\\"+"dilink.reg").Run()
				fmt.Println("Dilink installed for windows.")
				os.Exit(0)
			}
			case "darwin": {
				fmt.Println("need setting on web browser.")
				os.Exit(0)
			}
			case "linux": {
				exec.Command("gconftool-2", "--set", "/desktop/gnome/url-handlers/dilink/command", "--type=string", "/lustre/INHouse/CentOS/bin/dilink %s").Run()
				exec.Command("gconftool-2", "--set", "--type=bool", "/desktop/gnome/url-handlers/dilink/enabled", "true").Run()
				exec.Command("gconftool-2", "--set", "--type=bool", "/desktop/gnome/url-handlers/dilink/need-terminal", "false").Run()
				fmt.Println("Dilink installed for linux.")
				os.Exit(0)
			}
			default: {
				os.Exit(0)
			}
		}
	}

	//process protocol
	argstr = strings.Replace(flag.Args()[0], "dilink://","",1)
	switch runtime.GOOS {
		case "darwin": {
			if strings.HasSuffix(argstr, ".nk") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				os.Setenv("NUKE_PATH","/lustre/INHouse/nuke")
				os.Setenv("NUKE_OFX","/usr/OFX")
				exec.Command(RV_osx, scape).Run()
			} else if strings.HasSuffix(argstr, ".mov") || strings.HasSuffix(argstr, ".jpg") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command(RV_osx, scape).Run()
			} else if strings.HasSuffix(argstr, ".abc") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/lustre/INHouse/CentOS/bin/abcview", scape).Run()
			} else if strings.HasPrefix(argstr, "exec:") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("open", strings.Split(scape, ":")[1]).Run()
			} else {
				scape, err := url.QueryUnescape(argstr)
				fmt.Println(scape)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("open", scape).Run()
			}

		}
		case "linux": {
			if strings.HasSuffix(argstr, ".nk") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(scape, "lady") {
					os.Setenv("NUKE_PATH","/lustre/INHouse/nuke")
					os.Setenv("NUKE_OFX","/usr/OFX")
					os.Setenv("OPTICAL_FLARES_LICENSE_SERVER_IP","10.0.99.15")
					os.Setenv("BROWSER","firefox")
					os.Setenv("NUKE_FONT_PATH", "/lustre2/Digitalidea_source/2d_team_source/font")
					exec.Command("/usr/local/Nuke10.0v3/Nuke10.0", "--nukex", scape).Run()
				} else {
					os.Setenv("NUKE_PATH","/lustre/INHouse/nuke")
					os.Setenv("NUKE_OFX","/usr/OFX")
					os.Setenv("OPTICAL_FLARES_LICENSE_SERVER_IP","10.0.99.15")
					os.Setenv("BROWSER","firefox")
					os.Setenv("NUKE_FONT_PATH", "/lustre2/Digitalidea_source/2d_team_source/font")
					exec.Command("/usr/local/Nuke9.0v7/Nuke9.0", "--nukex", scape).Run()
				}

			} else if strings.HasSuffix(argstr, ".mov") || strings.HasSuffix(argstr, ".jpg") || strings.HasSuffix(argstr, ".png") || strings.HasSuffix(argstr, ".exr") || strings.HasSuffix(argstr, ".tga") || strings.HasSuffix(argstr, ".psd") || strings.HasSuffix(argstr, ".dpx") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(scape, "\\") && strings.Contains(scape, "\\\\10.0.200.100\\") {
					scape = strings.Replace(scape, "\\", "/", -1)
					scape = strings.Replace(scape, "//10.0.200.100/show_","/show/",1)
				}
				os.Setenv("RV_ENABLE_MIO_FFMPEG","1") // for prores

				if strings.Contains(scape, ";") {
					scapelist := strings.Split(scape, ";")

					exec.Command(RV_lin, scapelist...).Run()
				} else {
					exec.Command(RV_lin, scape).Run()
				}
			} else if strings.HasSuffix(argstr, ".rv") || strings.HasSuffix(argstr, ".mp4") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(scape, "\\") && strings.Contains(scape, "\\\\10.0.200.100\\") {
					scape = strings.Replace(scape, "\\", "/", -1)
					scape = strings.Replace(scape, "//10.0.200.100/show_","/show/",1)
				}
				os.Setenv("RV_ENABLE_MIO_FFMPEG","1") // for prores
				exec.Command(RV_lin, scape).Run()
			} else if strings.HasSuffix(argstr, ".ttf") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/usr/bin/gnome-font-viewer", scape).Run()

			} else if strings.HasSuffix(argstr, ".abc") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/lustre/INHouse/CentOS/bin/abcview", scape).Run()
			} else if strings.HasSuffix(argstr, ".py") || strings.HasSuffix(argstr, ".go") || strings.HasSuffix(argstr, ".txt") || strings.HasSuffix(argstr, ".md") || strings.HasSuffix(argstr, ".html") || strings.HasSuffix(argstr, ".css") || strings.HasSuffix(argstr, ".css") || strings.HasSuffix(argstr, ".env") || strings.HasSuffix(argstr, ".reg"){
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("gvim", "-R", scape).Run()

			} else if strings.HasSuffix(argstr, ".pdf") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/usr/bin/evince", scape).Run()
			} else if strings.HasSuffix(argstr, ".avi") || strings.HasSuffix(argstr, ".mkv"){
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/usr/bin/vlc", scape).Run()
			} else if strings.HasSuffix(argstr, ".blend") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/lustre/Applications/Linux/blender/blender-2.75a-linux-glibc211-x86_64/blender", scape).Run()

			} else if strings.HasSuffix(argstr, ".obj") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/lustre/Applications/Linux/blender/blender-2.75a-linux-glibc211-x86_64/blender", "--python", "/lustre/INHouse/blender/python/loadobj.py", "--", scape).Run()

			} else if strings.HasSuffix(argstr, ".hipnc") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/opt/hfs14.0.291/bin/happrentice", scape).Run()
			} else if strings.HasSuffix(argstr, ".hip") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/opt/hfs14.0.291/bin/houdinifx", scape).Run()

			} else if strings.HasSuffix(argstr, ".ntp") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("/lustre/Applications/Linux/natron/Natron2/Natron", scape).Run()
			} else if argstr == "bash" {
				exec.Command("gnome-terminal").Run()
			} else if strings.HasPrefix(argstr, "exec:") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("gnome-terminal", "-x", strings.Split(scape, ":")[1]).Run()
			} else {
				scape, err := url.QueryUnescape(argstr)
				fmt.Println(scape)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("nautilus", scape).Run()
			}
		}
		default: { //windows
			if strings.HasSuffix(argstr, ".mov") || strings.HasSuffix(argstr, ".rv") {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				if strings.Contains(scape, ";") {
					var movlist []string
					pathlist := strings.Split(scape, ";")
					for _, i := range pathlist {
						movlist = append(movlist, dipath.Lin2win(i))
					}
					exec.Command(RV_win, movlist...).Run()
				} else {
					exec.Command(RV_win, dipath.Lin2win(scape)).Run()
				}
			} else {
				scape, err := url.QueryUnescape(argstr)
				if err != nil {
					log.Fatal(err)
				}
				exec.Command("cmd","/C", "start", "", dipath.Lin2win(scape)).Run()
			}
		}
	}
}
