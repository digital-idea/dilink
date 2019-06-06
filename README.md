# dilink

![travisCI](https://secure.travis-ci.org/digital-idea/dilink.png)

웹에서 사용하기위한 웹 프로토콜입니다.
응용프로그램을 웹에서 실행하기 위해서 제작되었습니다.
윈도우즈에서는 admin 계정으로 등록해주세요.

#### Install
```
go get -u https://github.com/digital-idea/dipath.git
sh build.sh
```

#### 프로토콜 설치(CentOS6, Windows7)
```bash
$ dilink install
```

#### 프로토콜 설치(CentOS7)
터미널을 열고 아래처럼 명령어를 타이핑 합니다.
- 개발자

```bash
$ tcsh install_CentOS7_dev.sh
```

- 사용자

```bash
$ tcsh install_CentOS7.sh
```

#### HISTORY
- '17.8.7 : 입체조건이면 입체 프리뷰 지원하도록 기능추가.
- '15.4.23 : rvlink의 불편한 점들때문에 프로토콜 작성. 1차완료.
