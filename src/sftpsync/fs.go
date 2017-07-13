package sftpsync

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var localFileMap = make(map[string]*os.FileInfo)
var dataDir string

type LocalFs struct {
	dataDir string
	cache   map[string]*os.FileInfo
}

func (fs *LocalFs) Reload() error {
	for k, _ := range localFileMap {
		delete(localFileMap, k)
	}
	dir := fs.dataDir
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Log.Error(err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		Log.Debug("add cache file:", path)
		rel, err := filepath.Rel(dir, path)
		Log.Debug("rel file key:", rel)
		localFileMap[rel] = &info
		return nil
	})
}
func NewLocalFs(dir string) (*LocalFs, error) {
	dataDir = dir
	fs := LocalFs{dataDir: dir, cache: localFileMap}
	err := fs.Reload()
	//	for k, v := range localFileMap {
	//		fmt.Println(k, *v)
	//	}
	return &fs, err
}

func GetFileStat(name string) (stat os.FileInfo, err error) {
	if pStat, ok := localFileMap[name]; ok {
		return *pStat, nil
	}
	if ok, err := isFileExists(name); !ok {
		return nil, err
	}
	lStat, err := os.Lstat(getFileAbsPath(name))
	if err != nil {
		return
	}
	localFileMap[name] = &lStat
	return lStat, nil
}

func getFileAbsPath(name string) string {
	return filepath.Join(dataDir, name)
}

func isFileExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return true, nil
		}
		return false, errors.New(path + " exists but is not regular file")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkLocalDir(p string) error {
	fileInfo, err := os.Stat(p)
	if err != nil && !os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("stat dir:%s err:%v\n", p, err))
	}
	if err == nil {
		if fileInfo.IsDir() {
			return nil
		} else {
			os.Remove(p)
		}
	}
	if err := os.Mkdir(p, 0666); err != nil {
		return errors.New(fmt.Sprintf("create dir:%s err:%v\n", p, err))
	}
	return nil
}
func isFileChanged(rStat, lStat os.FileInfo) bool {
	if rStat == nil || lStat == nil {
		return false
	}
	if rStat.ModTime().Unix() == lStat.ModTime().Unix() {
		return true
	}
	return false
}
func saveLocalFile(reader io.Reader, p string) error {

	var perm os.FileMode = 0666
	file, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, perm)
	if err != nil {
		return errors.New(fmt.Sprintf("open local file:%s err:%v", p, err))
	}
	defer file.Close()
	b := make([]byte, 1024)
	for {
		n, err := reader.Read(b)
		file.Write(b[:n])
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
	return nil
}
