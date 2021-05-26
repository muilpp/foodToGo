test:
	go test -v ./...

compile:
	echo "Compiling for Linux, Mac and Win"
	GOOS=linux go build -o bin/food2GoLinux;
	GOOS=darwin GOARCH=amd64 go build -o bin/food2GoMac;
	GOOS=windows GOARCH=amd64 go build -o bin/food2GoWindows;

all: test compile