package internal

import (
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func TestNewSSHProxyHttpClient(t *testing.T) {
	user := "root"
	password := "123456"
	addr := "127.0.0.1:22"
	timeout := 15 * time.Second
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		Timeout:         timeout,
	}
	client, err := NewSSHTimeout("tcp", addr, config, timeout)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()
	httpClient := NewSSHProxyHttpClient(client)
	response, err := httpClient.Get("https://github.com/woluoyifan/go-ssh-special")
	if err != nil {
		t.Error(err)
		return
	}
	if response.StatusCode != 200 {
		t.Errorf("status code: %d", response.StatusCode)
	}
}
