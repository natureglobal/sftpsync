package sftpsync

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type sftpConfig struct {
	hostPort, user string
	password       string
	identityFile   string
}

func newSFTPClient(sc *sftpConfig) (*sftp.Client, error) {
	if sc.identityFile == "" && sc.password == "" {
		sc.password = prompter.Password(fmt.Sprintf("%s@%s's password", sc.user, sc.hostPort))
	}

	conf := &ssh.ClientConfig{
		User:            sc.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // XXX
	}
	if sc.identityFile != "" {
		out, err := os.ReadFile(sc.identityFile)
		if err != nil {
			return nil, err
		}
		signer, err := ssh.ParsePrivateKey(out)
		if err != nil {
			return nil, err
		}
		conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		conf.Auth = []ssh.AuthMethod{ssh.Password(sc.password)}
	}
	h := sc.hostPort
	if !strings.Contains(h, ":") {
		h += ":22"
	}
	sshConn, err := ssh.Dial("tcp", h, conf)
	if err != nil {
		return nil, err
	}
	return sftp.NewClient(sshConn)
}

func removeDir(cli *sftp.Client, d string) error {
	var dirs []string
	walker := cli.Walk(d)

	for walker.Step() {
		if err := walker.Err(); err != nil {
			return err
		}
		p := walker.Path()
		fi := walker.Stat()
		if fi.IsDir() {
			dirs = append(dirs, p)
			continue
		}
		log.Printf("removing: %s\n", p)
		if err := cli.Remove(p); err != nil {
			return err
		}
	}
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[j]) < len(dirs[i])
	})

	for _, d := range dirs {
		log.Printf("removing directory: %s\n", d)
		if err := cli.RemoveDirectory(d); err != nil {
			return err
		}
	}
	return nil
}
