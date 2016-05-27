.PHONY : build deps start-daemon

build:
	echo "Building Deploy It"
	go get && go build -o /opt/bin/deploy

deps:
	echo "Installing dependencies"
	go get

start-daemon:
	echo "Starting Deploy It daemon"
	go get && go build -o /opt/bin/deploy && deploy daemon --debug