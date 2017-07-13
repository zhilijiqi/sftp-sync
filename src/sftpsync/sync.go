package sftpsync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Syncd(conf *Config, req <-chan struct{}, done chan<- struct{}) {
	for {
		<-req
		sync(conf.Host, conf.User, conf.Password, conf.ServerPath, conf.LocalPath)
		done <- struct{}{}
	}
}

func sync(host, user, password, serverPath, localPath string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			Log.Error("sync err:%v", err)

		}
	}()

	client, err := getConn(user, password, host)
	if err != nil {
		Log.Error("Dial err:%v", err)
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		Log.Error("new session err:%v", err)
		return err
	}
	defer session.Close()

	ftpClient, err := sftp.NewClient(client)
	if err != nil {
		return err
	}
	defer ftpClient.Close()
	wd, _ := ftpClient.Getwd()
	Log.Debug("remote dir is:", wd, "remote sync dir is:", serverPath)

	walk := ftpClient.Walk(serverPath)
	for walk.Step() {
		if walk.Err() != nil {
			continue
		}
		if err := checkRemoteFile(ftpClient, walk.Stat(), walk.Path(), localPath); err != nil {
			Log.Error("remote file :%s err:%v\n", walk.Path(), err)
		}
	}
	return nil
}

func getConn(user, password, host string) (client *ssh.Client, err error) {
	authMethod := ssh.Password(password)
	auth := []ssh.AuthMethod{authMethod}
	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey()}

	client, err = ssh.Dial("tcp", host, conf)
	return
}

func checkRemoteFile(ftpClient *sftp.Client, rStat os.FileInfo, rPath, sftpSavePath string) error {
	p := filepath.FromSlash(sftpSavePath + rPath)
	if rStat.IsDir() {
		return checkLocalDir(p)
	}
	if ok, err := checkCacheFileChanged(rStat, rPath); !ok {

	}
	//	lStat, err := os.Lstat(p)
	//	if err != nil && !os.IsNotExist(err) {
	//		return errors.New(fmt.Sprintf("stat local file:%s err:%v", p, err))
	//	}
	//	if err == nil && isFileChanged(lStat, rStat) {
	//		return nil
	//	}
	sftpFile, err := ftpClient.Open(rPath)
	if err != nil {
		return errors.New(fmt.Sprintf("read sftp file:%s err:%v", p, err))
	}
	defer sftpFile.Close()
	if err := saveLocalFile(sftpFile, p); err != nil {
		return err
	}
	if err := os.Chtimes(p, time.Now(), rStat.ModTime()); err != nil {
		return errors.New(fmt.Sprintf("mtime file:%s err:%v", p, err))
	}
	return nil
}
