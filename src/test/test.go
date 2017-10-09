package main

import (
	"fmt"
	"runtime"
)

func main() {
	var int = runtime.NumGoroutine()
	for {
		fmt.Println(int)
	}
}
