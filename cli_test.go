package sftpsync_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/natureglobal/sftpsync"
)

func TestRun(t *testing.T) {
	hostPort := os.Getenv("TEST_SFTP_HOST")
	user := os.Getenv("TEST_SFTP_USER")
	key := os.Getenv("TEST_SFTP_KEY")

	err := sftpsync.Run(context.Background(),
		[]string{"-src=testdata/src", "-dst=upload/dst", fmt.Sprintf("-i=%s", key), user + "@" + hostPort},
		io.Discard, io.Discard)
	if err != nil {
		t.Error(err)
	}

	if err := sftpsync.Run(context.Background(), []string{"test@localhost"}, io.Discard, io.Discard); err == nil {
		t.Errorf("error shouldn't be nil")
	}

	if err := sftpsync.Run(context.Background(), []string{"-version"}, io.Discard, io.Discard); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
}
