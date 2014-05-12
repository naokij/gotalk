/*
Copyright 2014 Jiang Le

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package filestore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalProvider struct {
	Config *Config
}

var localpder = &LocalProvider{}

func (p *LocalProvider) Init(config Config) (err error) {
	if config.FSPath == "" {
		return fmt.Errorf("filestore:local: config.FSPath not provided")
	}
	//check if FSPath exists
	if file, err := os.Open(config.FSPath); os.IsNotExist(err) {
		return fmt.Errorf("filestore:local config.FSPath %s not exists", config.FSPath)
	} else {
		fileinfo, err2 := file.Stat()
		if err2 != nil {
			return fmt.Errorf("filestore:local unable to Stat config.FSPath %s", config.FSPath)
		}
		if fileinfo.IsDir() == false {
			return fmt.Errorf("filestore:local config.FSPath %s is not a directory", config.FSPath)
		}
		//is the path writeable?
	}
	fTest, err := os.Create(config.FSPath + "FSLocal.test")
	defer func() {
		fTest.Close()
		os.Remove(config.FSPath + "FSLocal.test")
	}()
	if err != nil {
		return fmt.Errorf("filestore:local config.FSPath %s is not writeable", config.FSPath)
	}
	p.Config = &config
	return nil
}

func (p *LocalProvider) PutFile(localFileUrl string, remoteFileUrl string) (url string, err error) {
	localFile, err := os.Open(localFileUrl)
	defer localFile.Close()
	if err != nil {
		return "", fmt.Errorf("filestore:local LocalFile %s not accessible", localFileUrl)
	}
	dir := filepath.Dir(remoteFileUrl)
	if dir != "." {
		if err := os.MkdirAll(p.Config.FSPath+dir, os.ModePerm); err != nil {
			return "", fmt.Errorf("filestore:local error while creating %s", p.Config.FSPath+dir)
		}
	}
	remoteFile, err2 := os.Create(p.Config.FSPath + remoteFileUrl)
	defer remoteFile.Close()
	if err2 != nil {
		return "", fmt.Errorf("filestore:local failed creating %s", p.Config.FSPath+remoteFileUrl)
	}
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return "", fmt.Errorf("filestore:local failed putting %s", p.Config.FSPath+remoteFileUrl)
	}
	return p.Config.UrlPrefix + remoteFileUrl, nil
}
func (p *LocalProvider) Get(filename string) (data []byte, err error) {
	return []byte(""), nil
}

func (p *LocalProvider) Delete(filename string) (err error) {
	err = os.Remove(p.Config.FSPath + filename)
	return err
}
func (p *LocalProvider) List(path string) (entries []Entry, err error) {
	return make([]Entry, 0), nil
}

func (p *LocalProvider) GetConfig() *Config {
	return p.Config
}

func init() {
	Register("local", localpder)
}
