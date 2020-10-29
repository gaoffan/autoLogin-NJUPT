default:
	GOOS=darwin GOARCH=amd64 go build -o out/autoLogin_darwin_amd64
	GOOS=linux GOARCH=amd64 go build -o out/autoLogin_linux_amd64
	GOARM=5 GOOS=linux GOARCH=arm go build -o out/autoLogin_linux_armv7
	GOOS=linux GOARCH=mips go build -o out/autoLogin_linux_mips
	GOOS=windows GOARCH=amd64 go build -o out/autoLogin_win_amd64.exe