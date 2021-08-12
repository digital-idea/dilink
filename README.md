# dilink

![travisCI](https://secure.travis-ci.org/digital-idea/dilink.svg)

웹에서 사용하기위한 웹 프로토콜입니다.
응용프로그램을 웹에서 실행하기 위해서 제작되었습니다.

### dilink 설치

#### Windows10

- `C:\bin` 폴더를 생성합니다.
- 위 폴더에 `dilink.exe`, `install_Windows.reg` 파일을 복사합니다.
- `install_Windows.reg` 파일을 더블클릭 합니다.

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

### Test실행
터미널을 이용해서 실제로 dilink가 잘 작동되는지 체크해볼 수 있습니다.

```bash
$ dilink dilink:///file/path/test.blend
```

### Mate Desktop 설치
MATE Desktop은 데스크탑 환경중 하나이다. CentOS7 데스크탑 환경이 GNOME3 으로 올라가면서 과거 GNOME2 환경이 그리운 사람들이 많이 사용한다.

```bash
$ sudo yum install -y epel-release
$ sudo yum groupinstall -y "MATE Desktop"
$ sudo reboot
```

### License
BSD 3-Clause License

### HISTORY
- '19.06.05: 오픈소스로 변경됨
- '17.08.07: 입체조건이면 RV에서 입체 프리뷰 지원하도록 기능 추가함.
- '15.04.23: rvlink의 불편한 점들때문에 dilink 프로토콜 작성.
