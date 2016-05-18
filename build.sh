GOOS=linux GOARCH=amd64 go build -o /lustre/INHouse/CentOS/bin/dilink dilink.go
GOOS=windows GOARCH=amd64 go build -o /lustre/INHouse/Windows/bin/dilink.exe dilink.go
GOOS=darwin GOARCH=amd64 go build -o /lustre/INHouse/OSX/bin/dilink dilink.go
