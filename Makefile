go_ldflags=-ldflags="-s -w"
agent_dist=dist
agent_main=main.go

build:
	rm -rf $(agent_dist)
	go mod tidy
	GOARCH=amd64 GOOS=linux go build -tags netcgo ${go_ldflags} -o $(agent_dist)/linux-amd64/agent_linux_amd64 $(agent_main)
	GOARCH=amd64 GOOS=windows go build -tags netcgo ${go_ldflags} -o $(agent_dist)/windows/agent_windows.exe $(agent_main)
	GOARCH=amd64 GOOS=darwin go build -tags netcgo ${go_ldflags} -o $(agent_dist)/mac-amd64/agent_mac_amd64 $(agent_main) 
	GOARCH=arm64 GOOS=darwin go build -tags netcgo ${go_ldflags} -o $(agent_dist)/mac-arm64/agent_mac_arm64 $(agent_main) 
	GOARCH=arm64 GOOS=linux go build -tags netcgo ${go_ldflags} -o $(agent_dist)/linux-arm64/agent_linux_arm64 $(agent_main)

compress_build:
	upx -9 $(agent_dist)/linux
	upx -9 $(agent_dist)/darwin
	upx -9 $(agent_dist)/windows.exe
	upx -9 $(agent_dist)/darwin-arm64

package: build
	tar zcvf dist/go-rest-agent.tgz dist/*
