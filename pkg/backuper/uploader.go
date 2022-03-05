package backuper

import (
	"bytes"
	"fmt"
	"os/exec"
	"time"
)

type uploader struct {
	aliOSS    *AlibabaOSS
	date      time.Time
	objKey    string
	cmdBin    string
	cmdArgs   []string
	maxRotate int
	env       []string
}

func (u *uploader) upload() error {
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

	err = u.aliOSS.Upload(u.objKey, gz)
	if err != nil {
		return fmt.Errorf("upload to alioss: %w", err)
	}

	return nil
}
