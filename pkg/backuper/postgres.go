package backuper

type PostgreCfg struct {
	User     string
	DBName   string
	Password string
}

type Postgre struct {
	cfg *PostgreCfg
}

func NewPostgre(cfg *PostgreCfg) *Postgre {
	return &Postgre{cfg: cfg}
}

const pgdumpBin = "pgdump"

// Backup backup postgres db to S3 compatible object storage
func (p *Postgre) Backup() error {

	return nil
}
