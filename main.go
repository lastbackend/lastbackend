package main

import (
	c "github.com/lastbackend/lastbackend/pkg/client/cmd"
	d "github.com/lastbackend/lastbackend/pkg/daemon/cmd"
)

func main() {
	d.Run()
	c.Run()
}
