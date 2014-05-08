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

type LocalProvider struct {
}

var localpder = &LocalProvider{}

func (p *LocalProvider) Init(config string) (err error) {
	return nil
}

func (p *LocalProvider) PutFile(localFileUrl string, remoteFileUrl string) (url string, err error) {
	return "", nil
}
func (p *LocalProvider) Get(filename string) (data []byte, err error) {
	return []byte(""), nil
}

func (p *LocalProvider) Delete(filename string) (files int, err error) {
	return 0, nil
}
func (p *LocalProvider) List(path string) (entries []Entry, err error) {
	return make([]Entry, 0), nil
}

func init() {
	Register("local", localpder)
}
