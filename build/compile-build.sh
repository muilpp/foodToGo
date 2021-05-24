cd ..
env GOOS=linux go build -o "build/food2GoLinux";
env GOOS=darwin GOARCH=amd64 go build -o "build/food2GoMac";
env GOOS=windows GOARCH=amd64 go build -o "build/food2GoWindows";
