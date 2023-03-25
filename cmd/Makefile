test:
	go test -v .././...

compile: test
	echo "Compiling for Linux and Rpi"
	GOOS=linux go build -o ../bin/food2GoLinux;
	GOOS=linux GOARCH=arm GOARM=5 go build -o ../bin/food2GoRpi;

deploy:
	scp /home/marc/workspace/foodToGo/bin/food2GoRpi rpi-mirall:/home/pi/foodToGo/

get-food:
	test compile
	pwd
	cd .. && ./bin/food2GoLinux getFood

all: test compile deploy
