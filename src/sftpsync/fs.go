package sftpsync

import (
	"container/list"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
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
	Log.Debug("FS Reload :", dir)

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
		rel = filepath.Join(string(filepath.Separator), rel)

		rel = filepath.FromSlash(rel)
		Log.Debug("rel file key:", rel)

		localFileMap[rel] = &info
		return nil
	})
}
func NewLocalFs(dir string) (*LocalFs, error) {
	dataDir = dir
	fs := LocalFs{dataDir: dir, cache: localFileMap}
	err := fs.Reload()
	for k, v := range localFileMap {
		Log.Debug("cache fs:", k, v)
	}
	return &fs, err
}

func checkCacheFileChanged(stat os.FileInfo, rPath string) (bool, error) {
	rPath = filepath.FromSlash(rPath)
	Log.Debug("checkCacheFileChanged path:", rPath)

	cacheStat, err := GetFileStat(rPath)
	if err != nil {
		return true, err
	}
	change := isFileChanged(cacheStat, stat)
	Log.Info("isFileChanged path:", change, rPath)
	return change, nil
}

func GetFileStat(name string) (stat os.FileInfo, err error) {
	Log.Debug("GetFileStat from cache:", name)
	if pStat, ok := localFileMap[name]; ok {
		Log.Info("GetFileStat from cache is Ok:", name)
		return *pStat, nil
	}
	absName := getLocalPath(name)
	Log.Debug("GetFileStat from path:", absName)
	if ok, err := isFileExists(absName); !ok {
		return nil, err
	}
	lStat, err := os.Lstat(absName)
	if err != nil {
		return
	}
	localFileMap[name] = &lStat
	return lStat, nil
}

func checkAndDelNotExistLocalFile(fileList *list.List) {
	if len(localFileMap) == fileList.Len() {
		return
	}
	waitDelFileMap := make(map[string]int)
	for e := fileList.Front(); e != nil; e = e.Next() {
		k := e.Value.(string)
		k = filepath.FromSlash(k)
		waitDelFileMap[k] = 0
	}
	for k, _ := range localFileMap {
		_, ok := waitDelFileMap[k]
		Log.Debug("waitDelFileMap", k)
		if !ok {
			delete(localFileMap, k)
			p := getLocalPath(k)
			os.Remove(p)
			Log.Info("del file:", p)
		}
	}
}

//func getFileAbsPath(name string) string {
//	return filepath.Join(dataDir, name)
//}

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
	p = getLocalPath(p)

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
		return true
	}
	if rStat.ModTime().Unix() == lStat.ModTime().Unix() {
		return false
	}
	return true
}

func getLocalPath(p string) string {
	return filepath.FromSlash(filepath.Join(dataDir, p))
}
func saveLocalFileAndMtime(reader io.Reader, p string, rStat os.FileInfo) error {
	p = getLocalPath(p)

	saveLocalFile(reader, p)
	if err := os.Chtimes(p, time.Now(), rStat.ModTime()); err != nil {
		return errors.New(fmt.Sprintf("mtime file:%s err:%v", p, err))
	}
	return nil
}
func saveLocalFile(reader io.Reader, p string) error {
	Log.Debug("saveLocalFile :", p)
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
