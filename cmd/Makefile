test:
	go test -v .././...

compile: test
	echo "Compiling for Linux, Mac and Win"
	GOOS=linux go build -o ../bin/food2GoLinux;
	GOOS=linux GOARCH=arm GOARM=5 go build -o ../bin/food2GoRpi;
	GOOS=darwin GOARCH=amd64 go build -o ../bin/food2GoMac;
	GOOS=windows GOARCH=amd64 go build -o ../bin/food2GoWindows;

deploy:
	scp /home/marc/workspace/foodToGo/bin/food2GoRpi rpi-mirall:/home/pi/foodToGo/

get-food:
	test compile
	pwd
	cd .. && ./bin/food2GoLinux getFood

all: test compile deploy
