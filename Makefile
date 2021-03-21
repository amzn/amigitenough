compile:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/amigitenough cmd/amigitenough/main.go
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/amigitenough cmd/amigitenough/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/amigitenough cmd/amigitenough/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/osx-amd64/amigitenough cmd/amigitenough/main.go
