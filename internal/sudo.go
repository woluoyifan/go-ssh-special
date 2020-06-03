package internal

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"regexp"
	"strings"
	"time"
)

func SudoExec(client *ssh.Client, password string, commands ...string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	session.Stdout = buf
	session.Stderr = buf
	in, _ := session.StdinPipe()
	go func(session *ssh.Session, in io.Writer, buf *bytes.Buffer) {
		r := regexp.MustCompile("\\[sudo\\] password for .+?: \r\nSorry, try again.\r\n")
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			s := buf.String()
			if strings.Contains(s, "[sudo] password for ") {
				_, err = in.Write([]byte(password + "\n"))
				//password err check
				for i := 0; i < 5; i++ {
					if r.MatchString(s) {
						_ = session.Close()
						break
					}
				}
				return
			}
		}
	}(session, in, buf)
	command := fmt.Sprintf("sudo -S %s", strings.Join(commands, "\n"))
	err = session.Run(command)
	if err != nil {
		return "", err
	}
	s := buf.String()
	i := strings.Index(s, "\n")
	return s[i+1:], nil
}
