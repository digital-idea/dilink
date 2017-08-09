package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 참고명령 rv [ /show/mkk3/seq/BNS/BNS_0260/plate/BNS_0260_left.mov /show/mkk3/seq/BNS/BNS_0260/plate/BNS_0260_right.mov ] [ /show/mkk3/seq/BNS/BNS_0200/plate/BNS_0200_left.mov /show/mkk3/seq/BNS/BNS_0200/plate/BNS_0200_right.mov ] -stereo scanline
// 만약 짝이 있다면
// "/show/test_left.mov" -> "[ /show/test_left.mov /show/test_right.mov ]" 로 리턴되어야 한다.
// 만약 짝이 없다면 그대로 리턴된다.
// "/show/test_left.mov" -> "/show/test_left.mov" 로 리턴된다.
// 입체가 1개라도 섞여있다면(짝) rv실행시 "-stereo scanline" 인수를 마지막 리스트에 넣는다.

// RV의 mov리스트 인수값을 받아서 입체셋팅으로 재생할 상황이면 입체 재생을 위한 인수 리스트를 반환한다.
func ToRvStereo(movlist []string) []string {
	// 만약 빈 리스트라면 바로 리턴한다.
	if len(movlist) == 0 {
		return movlist
	}
	playlist := []string{}
	isStereo := false
	for _, movfile := range movlist {
		dir, file := filepath.Split(movfile)
		// mov파일이 존재하는가?
		if _, err := os.Stat(movfile); os.IsNotExist(err) {
			continue
		}

		// 입체 프로젝트시 회사 기준인 "left" 문자가 파일명에 존재하는가?
		if strings.Contains(movfile, "left") {
			rightmov := filepath.Join(dir, strings.Replace(file, "left", "right", -1))
			if _, err := os.Stat(rightmov); os.IsNotExist(err) {
				// 입체이지만 left 영상만 렌더링 걸려있는 상황이다. left 영상이지만 일반영상처럼 취급한다.
				playlist = append(playlist, movfile)
				continue
			}
			playlist = append(playlist, "[")
			playlist = append(playlist, movfile)
			playlist = append(playlist, rightmov)
			playlist = append(playlist, "]")
			// 입체조건이 성립되었다.
			if !isStereo {
				isStereo = true
			}
			continue
		}

		// left 규약을 어기고 작업자가 "right" 문자를 CSI에 등록했을때의 상황도 대비한다.
		if strings.Contains(movfile, "right") {
			leftmov := filepath.Join(dir, strings.Replace(file, "right", "left", -1))
			if _, err := os.Stat(leftmov); os.IsNotExist(err) {
				// 입체이지만 right 영상만 렌더링 걸려있는 상황이다. right 영상이지만 일반영상처럼 취급한다.
				playlist = append(playlist, movfile)
				continue
			}
			playlist = append(playlist, "[")
			playlist = append(playlist, leftmov)
			playlist = append(playlist, movfile)
			playlist = append(playlist, "]")
			// 입체조건이 성립되었다.
			if !isStereo {
				isStereo = true
			}
		}
	}

	if isStereo {
		// 입체 컨펌시 RV에서 "-stereo scanline" 옵션이  입체모니터 컨펌에 맞다.
		playlist = append(playlist, "-stereo")
		playlist = append(playlist, "scanline")
	}
	fmt.Println(playlist)
	return playlist
}
