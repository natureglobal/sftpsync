package sftpsync

import (
	"os"
	"testing"
)

func TestSFTP(t *testing.T) {
	sc := &sftpConfig{
		hostPort: os.Getenv("TEST_SFTP_HOST"),
		user:     os.Getenv("TEST_SFTP_USER"),
		password: os.Getenv("TEST_SFTP_PASSWORD"),
	}
	_, err := newSFTPClient(sc)
	if err != nil {
		t.Error(err)
	}
}
