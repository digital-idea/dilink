# dilink

![travisCI](https://secure.travis-ci.org/digital-idea/dilink.png)

웹에서 사용하기위한 웹 프로토콜입니다.
응용프로그램을 웹에서 실행하기 위해서 제작되었습니다.

### Download
- [Linux 64bit](https://github.com/digital-idea/dilink/releases/download/v1.0.3/dilink_linux_x86-64.tgz)
- [macOS 64bit](https://github.com/digital-idea/dilink/releases/download/v1.0.3/dilink_darwin_x86-64.tgz)
- [Windows 64bit](https://github.com/digital-idea/dilink/releases/download/v1.0.3/dilink_windows_x86-64.zip)

### dilink 설치

#### Windows10

`C:\bin\dilink.exe` 에 파일을 복사합니다.
이후 install_Windows.reg 레지스트리를 실행합니다.

#### CentOS7
터미널을 열고 아래처럼 명령어를 타이핑 합니다.

```bash
$ tcsh install_CentOS7.sh // 사용자
$ tcsh install_CentOS7_dev.sh // 개발자
```

#### macOS
1. 다운로드 받은 파일을 압축풉니다.
1. `dilink.app` 파일을 어플리케이션에 복사합니다.
1. dilink 명령어는 ~/bin 폴더에 넣습니다.


### License
BSD 3-Clause License

### HISTORY
- '19.06.05: 오픈소스로 변경됨
- '17.08.07: 입체조건이면 RV에서 입체 프리뷰 지원하도록 기능 추가함.
- '15.04.23: rvlink의 불편한 점들때문에 dilink 프로토콜 작성.
