package backuper

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/rs/zerolog/log"
)

const mysqldumpBin = "mysqldump"
const dateLayout = "2006-01-02"

type MySQLCfg struct {
	User     string
	DBName   string
	Password string
	OutDir   string
}

type MySQL struct {
	cfg *MySQLCfg
}

func NewMySQL(cfg *MySQLCfg) *MySQL {
	return &MySQL{cfg: cfg}
}

func (m *MySQL) command() []string {
	cmd := fmt.Sprintf(`mysqldump -u %s -p%s --databases %s`, m.cfg.User, m.cfg.Password, m.cfg.DBName)
	return strings.Split(cmd, " ")
}

// Backup backup to git repo
func (m *MySQL) Backup() (outpath string, err error) {
	date := time.Now().Format(dateLayout)
	filename := fmt.Sprintf("%s-%s.postgres.gz", slug.Make(m.cfg.DBName), date)
	outpath = path.Join(m.cfg.OutDir, filename)

	cmdArgs := m.command()
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debug().Str("cmd", cmd.String()).Msg("")
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("run cmd: %s: %w", stderr.String(), err)
		return
	}

	err = writeBackup(stdout, outpath)
	if err != nil {
		return "", fmt.Errorf("save dump to file: %w", err)
	}

	return
}

// DockerBackup expect bash shell
func (m *MySQL) DockerBackup(containerID string) (outpath string, err error) {
	date := time.Now().Format(dateLayout)
	filename := fmt.Sprintf("%s-%s.mysql.sql.gz", slug.Make(m.cfg.DBName), date)
	outpath = path.Join(m.cfg.OutDir, filename)

	cmd := exec.Command("docker", "exec", "-i", containerID)
	cmd.Args = append(cmd.Args, m.command()...)

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	log.Debug().Str("cmd", cmd.String()).Msg("")
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%s: %w", stderr.String(), err)
		return
	}

	err = writeBackup(stdout, outpath)
	if err != nil {
		return "", fmt.Errorf("save dump to file: %w", err)
	}

	return
}

func writeBackup(src io.Reader, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer outFile.Close()

	gz, err := gzip.NewWriterLevel(outFile, gzip.DefaultCompression)
	if err != nil {
		return fmt.Errorf("create gzip writer: %w", err)
	}
	defer gz.Close()
	defer gz.Flush()

	_, err = io.Copy(gz, src)
	if err != nil {
		return fmt.Errorf("save backup to file: %w", err)
	}

	return nil
}
