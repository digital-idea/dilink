# dilink
- 웹에서 사용하기위한 dilink 프로토콜입니다.
- 응용프로그램, 기능을 웹에서 구현하기 위해서 제작되었습니다.
- 윈도우즈에서는 admin 계정으로 등록해주세요.

#### 기본정보
- 개발일시 : '15.12
- 관리부서 : 파이프라인팀

#### install
```
go get -u https://github.com/didev/dipath.git
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
