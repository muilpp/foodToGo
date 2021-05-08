env GOOS=linux go build -o "food2GoLinux";
env GOOS=darwin GOARCH=amd64 go build -o "food2GoMac";
env GOOS=windows GOARCH=amd64 go build -o "food2GoWindows";
