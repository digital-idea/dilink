#!/bin/sh
APP="dilink"

# OS별 필드
GOOS=windows GOARCH=amd64 go build -o ./bin/windows/${APP}.exe dilink.go stereo.go windows.go pathfunc.go

# Github Release에 업로드 하기위해 압축
cd ./bin/windows/ && cp ../../../../install/windows/install_Windows.reg . && zip -r ../${APP}_windows_x86-64.zip . && cd -

# 삭제
rm -rf ./bin/windows
