package backuper

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/gosimple/slug"
)

type DBCfg struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
	OutDir   string
}

type PostgreCfg struct {
	DBCfg
	AliOSS *AlibabaOSS
}

type Postgre struct {
	cfg *PostgreCfg
}

func NewPostgre(cfg *PostgreCfg) *Postgre {
	return &Postgre{cfg: cfg}
}

func (p *Postgre) command() (cmdBin string, args []string) {
	cmd := fmt.Sprintf(`pg_dump --dbname %s -w`, p.cfg.DBName)
	cmds := strings.Split(cmd, " ")
	return cmds[0], cmds[1:]
}

// Backup backup postgres db to S3 compatible object storage
func (p *Postgre) Backup() (outpath string, err error) {
	cmdBin, args := p.command()
	date := time.Now().Format(dateLayout)
	filename := fmt.Sprintf("%s-%s.pg.sql.gz", slug.Make(p.cfg.DBName), date)
	outpath = path.Join(p.cfg.OutDir, filename)

	cmd := exec.Command(cmdBin, args...)
	cmd.Env = append(cmd.Env,
		"PGHOST="+p.cfg.Host,
		"PGPORT="+fmt.Sprint(p.cfg.Port),
		"PGUSER="+p.cfg.User,
		"PGPASSWORD="+p.cfg.Password,
	)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("run cmd: %s: %w", stderr.String(), err)
		return
	}

	err = writeBackup(stdout, outpath)
	if err != nil {
		return "", fmt.Errorf("save dump to file: %w", err)
	}

	err = p.cfg.AliOSS.UploadFromPath(outpath, filename)
	if err != nil {
		return "", fmt.Errorf("postgre backup: upload to alioss: %w", err)
	}

	return
}
