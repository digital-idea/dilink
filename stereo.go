package main

import (
	"os"
	"path/filepath"
	"strings"
)

// RV에서 입체를 보기위한 테스트파일 : test/view_stereo.sh

// ToRvStereo함수는 입체조건을 체크하고,
// 조건이 맞다면 RV에서 입체로 재생할 수 있는 문자를 반환한다.
// 만약 /show/test_left.mov파일에 짝이 있다면
// []string{"[", "/show/test_left.mov", "/show/test_right.mov", "]"}, true 를 반환한다.
// 만약 /show/test_left.mov 파일에 짝이 없다면
// []string{}, false 를 반환한다.
func ToRvStereo(movfile string) ([]string, bool) {
	// mov파일이 존재하는가?
	if _, err := os.Stat(movfile); os.IsNotExist(err) {
		// 존재하지 않으면 movfile을 리턴한다.
		// 플레이 리스트에 다른 mov파일도 많이 있기 때문에서
		// RV에서 에러를 처리해야 플레이리스트 갯수가 깨지지 않는다.
		return []string{}, false
	}
	if !(strings.Contains(movfile, "left") || strings.Contains(movfile, "right")) {
		// 입체관련된 문자열이 파일명에 포함되어있지 않다.
		// 바로 반환한다.
		return []string{}, false
	}

	// 입체 프로젝트시 회사 기준인 "left" 문자가 파일명에 존재하는가?
	dir, file := filepath.Split(movfile)
	if strings.Contains(movfile, "left") {
		rightmov := filepath.Join(dir, strings.Replace(file, "left", "right", -1))
		if _, err := os.Stat(rightmov); os.IsNotExist(err) {
			// 입체이지만 left 영상만 렌더링 걸려있는 상황이다. left 영상이지만 일반영상처럼 취급한다.
			return []string{}, false
		}
		// 입체조건이 성립되었다.
		return []string{"[", movfile, rightmov, "]"}, true
	}

	// left 규약을 어기고 작업자가 "right" 문자를 CSI에 등록했을때의 상황도 대비한다.
	if strings.Contains(movfile, "right") {
		leftmov := filepath.Join(dir, strings.Replace(file, "right", "left", -1))
		if _, err := os.Stat(leftmov); os.IsNotExist(err) {
			// 입체이지만 right 영상만 렌더링 걸려있는 상황이다. right 영상이지만 일반영상처럼 취급한다.
			return []string{}, false
		}
		// 입체조건이 성립되었다.
		return []string{"[", leftmov, movfile, "]"}, true
	}
	return []string{}, false
}
