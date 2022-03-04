package model

type Backupper interface {
	Backup() (outpath string, err error)
}
