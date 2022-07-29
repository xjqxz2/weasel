.PHONY:local

# 本地环境编译
local:
	go build  -o $(shell pwd)/bin/debug $(shell pwd)
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(shell pwd)/bin/linux_x64-86 $(shell pwd)
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -buildmode=c-shared -ldflags "-s -w" -o $(shell pwd)/bin/windows_x64-86.exe $(shell pwd)
