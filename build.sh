CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /lustre/INHouse/CentOS/bin/dilink dilink.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /lustre/INHouse/Windows/bin/dilink.exe dilink.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o /lustre/INHouse/OSX/bin/dilink dilink.go
