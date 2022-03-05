package backuper

import (
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

type MySQLCfg struct {
	DBCfg
	AliOSS *AlibabaOSS
}

type MySQL struct {
	cfg *MySQLCfg
}

func NewMySQL(cfg *MySQLCfg) *MySQL {
	return &MySQL{cfg: cfg}
}

func (m *MySQL) command() (cmdBin string, args []string) {
	cmd := fmt.Sprintf(
		`mysqldump --host %s --port %d --protocol tcp --skip-column-statistics -u %s -p%s --databases %s`,
		m.cfg.Host,
		m.cfg.Port,
		m.cfg.User,
		m.cfg.Password,
		m.cfg.DBName)
	cmds := strings.Split(cmd, " ")
	return cmds[0], cmds[1:]
}

// Backup backup to git repo
func (m *MySQL) Backup() (outpath string, err error) {
	date := time.Now()
	ext := "mysql.sql.gz"
	filename := fmt.Sprintf("%s.%s", m.cfg.DBName, ext)
	rotateTag := makeRotateTag(m.cfg.MaxRotate, date, filename)
	cmdBin, cmdArgs := m.command()

	up := uploader{
		aliOSS:    m.cfg.AliOSS,
		date:      date,
		objKey:    rotateTag,
		cmdBin:    cmdBin,
		cmdArgs:   cmdArgs,
		maxRotate: m.cfg.MaxRotate,
	}

	return rotateTag, up.upload()
}
