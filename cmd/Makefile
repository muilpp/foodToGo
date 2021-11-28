test:
	go test -v .././...

compile:
	echo "Compiling for Linux, Mac and Win"
	GOOS=linux go build -o ../bin/food2GoLinux;
	GOOS=darwin GOARCH=amd64 go build -o ../bin/food2GoMac;
	GOOS=windows GOARCH=amd64 go build -o ../bin/food2GoWindows;

deploy:
	scp /home/marc/workspace/marc/get-food-to-go/bin/food2GoLinux marc-digital-ocean:/home/marc/foodToGo/

get-food:
	test compile
	pwd
	cd .. && ./bin/food2GoLinux getFood

all: test compile