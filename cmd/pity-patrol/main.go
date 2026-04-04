package main

import "atomicptr.dev/pity-patrol/pkg/cli"

func main() {
	err := cli.Run()
	if err != nil {
		panic(err)
	}
}
