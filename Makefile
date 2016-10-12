.PHONY : build deps

build:
	echo "Building Deploy It"
	go get && go build -o bin/deploy
	mkdir -p build/linux  && GOOS=linux  go build -ldflags "-X main.Version=$(VERSION)" -o build/linux/$(NAME)
	mkdir -p build/darwin && GOOS=darwin go build -ldflags "-X main.Version=$(VERSION)" -o build/darwin/$(NAME)

deps:
	echo "Installing dependencies"
	go get
