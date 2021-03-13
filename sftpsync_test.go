package sftpsync

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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
	ap := &app{
		url:      fmt.Sprintf("sftp://%s@%s", user, hostPort),
		password: password,
		src:      "testdata/src",
		dst:      "upload/dst",
	}
	if err := ap.run(context.Background(), ioutil.Discard, ioutil.Discard); err != nil {
		t.Error(err)
	}
}
