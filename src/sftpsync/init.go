package sftpsync

import (
	"fmt"
	"os"
	"path/filepath"
)

func init() {
	fmt.Println("init.....")
}

var baseDir = filepath.Dir(os.Args[0])
