package main

import "github.com/atomicptr/pity-patrol/pkg/cli"

func main() {
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
