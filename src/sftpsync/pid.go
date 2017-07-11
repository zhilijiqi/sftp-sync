package sftpsync

import (
	"os"
	"path/filepath"
	"strconv"
)

func GetPid() int {
	pid := os.Getpid()
	f, err := os.Create(filepath.Join(baseDir, "pid"))
	if err != nil {
		return -1
	}
	f.WriteString(strconv.Itoa(pid))
	defer f.Close()
	return pid
}
