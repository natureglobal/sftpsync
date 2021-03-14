package sftpsync

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestAppGetClient(t *testing.T) {
	hostPort := os.Getenv("TEST_SFTP_HOST")
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		t.Fatal(err)
	}
	user := os.Getenv("TEST_SFTP_USER")
	password := os.Getenv("TEST_SFTP_PASSWORD")

	testCases := []struct {
		name string
		app  *app
	}{{
		name: "scp like URL",
		app: &app{
			url:      fmt.Sprintf("%s@%s", user, host),
			password: password,
		},
	}, {
		name: "scp like URL with password",
		app: &app{
			url: fmt.Sprintf("%s:%s@%s", user, password, host),
		},
	}, {
		name: "sftp scheme",
		app: &app{
			url:      fmt.Sprintf("sftp://%s@%s", user, host),
			password: password,
		},
	}, {
		name: "sftp scheme with password and port",
		app: &app{
			url: fmt.Sprintf("sftp://%s:%s@%s:%s", user, password, host, port),
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := tc.app.getClient(context.Background()); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestAppGetRun(t *testing.T) {
	hostPort := os.Getenv("TEST_SFTP_HOST")
	user := os.Getenv("TEST_SFTP_USER")
	password := os.Getenv("TEST_SFTP_PASSWORD")

	dst := "upload/dst"
	ap := &app{
		url:      fmt.Sprintf("sftp://%s@%s", user, hostPort),
		password: password,
		src:      "testdata/src",
		dst:      dst,
	}
	cli, _, err := ap.getClient(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if err := ap.run(context.Background(), ioutil.Discard, ioutil.Discard); err != nil {
		t.Error(err)
	}
	for _, f := range []string{"file1", "dir1/file2"} {
		fpath := path.Join(dst, f)
		fi, err := cli.Stat(fpath)
		if err != nil || fi.Size() <= 0 {
			t.Errorf("something went wrong: %s", err)
		}
	}

	// overwrite
	if err := ap.run(context.Background(), ioutil.Discard, ioutil.Discard); err != nil {
		t.Error(err)
	}
	for _, f := range []string{"file1", "dir1/file2"} {
		fpath := path.Join(dst, f)
		fi, err := cli.Stat(fpath)
		if err != nil || fi.Size() <= 0 {
			t.Errorf("something went wrong: %s", err)
		}
	}

	_, err = cli.Stat(filepath.Join(dst, "src"))
	if err == nil {
		t.Errorf("something went wrong")
	}

	if err := removeDir(cli, dst); err != nil {
		t.Error(err)
	}

	_, err = cli.Stat(dst)
	if err == nil {
		t.Errorf("something went wrong")
	}
}
