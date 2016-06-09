.PHONY : build deps start-daemon

build:
	echo "Building Deploy It"
	go get && go build -o /opt/bin/deploy

deps:
	echo "Installing dependencies"
	go get