package main

import (
	"fmt"
	"log"

	ic "github.com/takanoriyanagitani/go-arrow-ipc-count"
)

func sub() error {
	cnt, e := ic.StdinToIpcNamesToTotalCountDefault()
	if nil != e {
		return e
	}

	fmt.Printf("%v\n", cnt) //nolint:forbidigo
	return nil
}

func main() {
	e := sub()
	if nil != e {
		log.Fatalf("%v\n", e)
	}
}
