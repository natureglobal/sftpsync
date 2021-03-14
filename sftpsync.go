package sftpsync

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/sftp"
)

type app struct {
	url          string
	port         uint
	password     string
	identityFile string
	src, dst     string
}

const cmdName = "sftpsync"

// Run the sftpsync
func Run(ctx context.Context, argv []string, outStream, errStream io.Writer) error {
	log.SetOutput(errStream)
	fs := flag.NewFlagSet(
		fmt.Sprintf("%s (v%s rev:%s)", cmdName, version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)

	ap := &app{password: os.Getenv("SFTP_PASSWORD")}
	ver := fs.Bool("version", false, "display version")
	fs.UintVar(&ap.port, "P", 0, "port")
	fs.StringVar(&ap.src, "src", "", "source directory")
	fs.StringVar(&ap.dst, "dst", "", "destination directory")
	fs.StringVar(&ap.identityFile, "i", "", "identity file")
	if err := fs.Parse(argv); err != nil {
		return err
	}
	if *ver {
		return printVersion(outStream)
	}
	ap.url = fs.Arg(0)

	return ap.run(ctx, outStream, errStream)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "%s v%s (rev:%s)\n", cmdName, version, revision)
	return err
}

func (ap *app) run(ctx context.Context, outStream, errStream io.Writer) error {
	cl, _, err := ap.getClient(ctx)
	if err != nil {
		return err
	}
	return filepath.WalkDir(ap.src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(ap.src, p)
		if err != nil {
			return err
		}

		dstPath := path.Join(ap.dst, relPath)
		if d.IsDir() {
			return cl.MkdirAll(dstPath)
		}
		log.Printf("transfer %s to %s\n", p, dstPath)
		src, err := os.Open(p)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := cl.Create(dstPath) // XXX not atomic
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		return err
	})
}

var scpURLReg = regexp.MustCompile("^([^@]+@)?([^:]+)(?::(/?.+))?$")

func (ap *app) getClient(ctx context.Context) (*sftp.Client, string, error) {
	if !strings.HasPrefix(ap.url, "sftp://") {
		if m := scpURLReg.FindStringSubmatch(ap.url); len(m) > 3 {
			ap.url = fmt.Sprintf("sftp://%s%s/%s", m[1], m[2], m[3])
		} else {
			ap.url = "sftp://" + ap.url
		}
	}
	u, err := url.Parse(ap.url)
	if err != nil {
		return nil, "", err
	}

	hostPort := u.Host
	if !strings.Contains(hostPort, ":") && ap.port > 0 {
		hostPort += fmt.Sprintf(":%d", ap.port)
	}
	pass, _ := u.User.Password()
	if pass == "" {
		pass = ap.password
	}

	sc := &sftpConfig{
		hostPort:     hostPort,
		user:         u.User.Username(),
		password:     pass,
		identityFile: ap.identityFile,
	}
	cli, err := newSFTPClient(sc)
	return cli, u.Path, err
}
