package filestore

import (
	"testing"
)

var localFS Manager

func init() {
}

func TestLocalFSInit(t *testing.T) {
	config := Config{UrlPrefix: "http://127.0.0.1:8080/avatars/", FSPath: "/tmp/"}
	localFS, err := NewManager("local", config)
	if err != nil {
		t.Error(err)
	}
	t.Log("filestore local: localfs init %s", localFS)
}

func TestLocalFSPut(t *testing.T) {
	config := Config{UrlPrefix: "http://127.0.0.1:8080/avatars/", FSPath: "/tmp/"}
	localFS, err := NewManager("local", config)
	if err != nil {
		t.Error(err)
	}
	var url string
	url, err = localFS.PutFile("/tmp/test.jpg", "remotefile.jpg")
	if err != nil {
		t.Error(err)
	}
	t.Log("file put successful url:" + url)
}
