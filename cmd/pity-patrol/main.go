package main

import "github.com/atomicptr/pity-patrol/pkgs/cli"

func main() {
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
