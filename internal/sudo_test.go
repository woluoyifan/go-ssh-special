package internal

import (
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func TestSudoExec(t *testing.T) {
	user := "test"
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
	result, err := SudoExec(client, password, "ls /root")
	if err !=nil{
		t.Error(err)
		return
	}
	println(result)
}
