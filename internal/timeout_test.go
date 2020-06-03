package internal

import (
	"errors"
	"golang.org/x/crypto/ssh"
	"testing"
	"time"
)

func TestNewSSHTimeout(t *testing.T) {
	user := "root"
	password := "123456"
	addr := "127.0.0.1:22"
	timeout := 15 * time.Second
	config := &ssh.ClientConfig{
		User: user,
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
	t.Log("try session1 ... ")
	session, err := client.NewSession()
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()
	output, err := session.CombinedOutput("ls")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("session1 result: %s", string(output))
	t.Logf("wait %v for timeout",timeout)
	<-time.After(timeout)
	t.Log("try session2 ... ")
	_, err = client.NewSession()
	if err == nil {
		t.Error(errors.New("client not timeout"))
		return
	}
	t.Log("session2 timeout pass")
}
