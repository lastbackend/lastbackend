.PHONY : build deps

build:
	echo "Building Deploy It"
	go get && go build -o /usr/local/bin/deploy

deps:
	echo "Installing dependencies"
	go get