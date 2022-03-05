package backuper

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

type ObjectStorageService interface {
	Upload(key string, rd io.Reader) error
}

type backuper struct {
	oss       ObjectStorageService
	date      time.Time
	objKey    string
	cmdBin    string
	cmdArgs   []string
	maxRotate int
	env       []string
}

// run backup command and upload to object storage service
func (u *backuper) backup() error {
	cmd := exec.Command(u.cmdBin, u.cmdArgs...)
	dbDump := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cmd.Stdout = dbDump
	cmd.Stderr = stderr
	cmd.Env = append(cmd.Env, u.env...)

	if err := cmd.Run(); err != nil {
		err = fmt.Errorf("run cmd: %s: %w", stderr.String(), err)
		return err
	}

	gz := bytes.NewBuffer(nil)
	err := writeGzip(dbDump, gz)
	if err != nil {
		return fmt.Errorf("gzip dump: %w", err)
	}

	err = u.oss.Upload(u.objKey, gz)
	if err != nil {
		return fmt.Errorf("upload to alioss: %w", err)
	}

	return nil
}
