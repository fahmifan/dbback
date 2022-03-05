package backuper

import (
	"fmt"
	"strings"
	"time"
)

const defaultMaxRotate = 7

type DBCfg struct {
	Host      string
	Port      int
	User      string
	DBName    string
	Password  string
	OutDir    string
	MaxRotate int
}

type PostgreCfg struct {
	DBCfg
	AliOSS *AlibabaOSS
}

type Postgre struct {
	cfg *PostgreCfg
}

func NewPostgre(cfg *PostgreCfg) *Postgre {
	if cfg.MaxRotate <= 0 {
		cfg.MaxRotate = defaultMaxRotate
	}
	return &Postgre{cfg: cfg}
}

func (p *Postgre) command() (cmdBin string, args []string) {
	cmd := fmt.Sprintf(`pg_dump --dbname %s -w`, p.cfg.DBName)
	cmds := strings.Split(cmd, " ")
	return cmds[0], cmds[1:]
}

// Backup backup postgres db to S3 compatible object storage
func (p *Postgre) Backup() (_ string, err error) {
	date := time.Now()
	filename := p.cfg.DBName + ".pg.sql.gz"
	rotateTag := makeRotateTag(p.cfg.MaxRotate, date, filename)
	env := []string{
		"PGHOST=" + p.cfg.Host,
		"PGPORT=" + fmt.Sprint(p.cfg.Port),
		"PGUSER=" + p.cfg.User,
		"PGPASSWORD=" + p.cfg.Password,
	}
	cmdBin, cmdArgs := p.command()

	up := backuper{
		oss:       p.cfg.AliOSS,
		date:      date,
		objKey:    rotateTag,
		cmdBin:    cmdBin,
		cmdArgs:   cmdArgs,
		maxRotate: p.cfg.MaxRotate,
		env:       env,
	}

	return rotateTag, up.backup()
}
