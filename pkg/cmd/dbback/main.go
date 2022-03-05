package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/fahmifan/dbback/pkg/backuper"
	"github.com/fahmifan/dbback/pkg/config"
	"github.com/fahmifan/dbback/pkg/model"
	"github.com/fahmifan/dbback/pkg/worker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: zerolog.TimeFieldFormat,
	}).With().Timestamp().Caller().Logger()

	if err := run(os.Args); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}

const menu = `
backupper is cli to backup db, currently support mysql & postgres

backupper [options]

options:
	--help		show the help menu
	--dbaname	the database name
	--driver 	the database driver [mysql, postgres]
	--cron		run as cron job
`

func run(args []string) error {
	if len(args) <= 1 || (len(args) > 1 && args[1] == "--help") {
		fmt.Print(menu)
		return nil
	}

	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		bak model.Backupper

		outputPath string
		dbName     string
		dbDriver   string
		isCron     bool
	)

	cmd.StringVar(&dbName, "dbname", "", `--dbname foobar`)
	cmd.StringVar(&dbDriver, "driver", "", `--driver [mysql, c, postgres]`)
	cmd.BoolVar(&isCron, "cron", false, `--cron [true, false] default false`)

	if err := cmd.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	cfg, err := config.Load("./config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ossClient, err := oss.New(cfg.AliOSS.Endpoint, cfg.AliOSS.AccessKeyID, cfg.AliOSS.AccessKeySecret, oss.Timeout(10, 30))
	if err != nil {
		return fmt.Errorf("new oss: %w", err)
	}

	aliOss := backuper.NewAlibabaOSS(&backuper.AlibabaOSSCfg{
		Client: ossClient,
		Bucket: cfg.AliOSS.Bucket,
	})

	switch dbDriver {
	default:
		return errors.New("invalid driver, should be [mysql, postgres]")
	case "postgres":
		bak = backuper.NewPostgre(&backuper.PostgreCfg{
			AliOSS: aliOss,
			DBCfg: backuper.DBCfg{
				OutDir:   cfg.OutDir,
				User:     cfg.Postgres.User,
				Password: cfg.Postgres.Password,
				Host:     cfg.Postgres.Host,
				Port:     cfg.Postgres.Port,
				DBName:   dbName,
			},
		})
	case "mysql":
		bak = backuper.NewMySQL(&backuper.MySQLCfg{
			AliOSS: aliOss,
			DBCfg: backuper.DBCfg{
				OutDir:   cfg.OutDir,
				User:     cfg.MySQL.User,
				Password: cfg.MySQL.Password,
				Host:     cfg.MySQL.Host,
				Port:     cfg.MySQL.Port,
				DBName:   dbName,
			},
		})
	}

	if !isCron {
		outputPath, err = bak.Backup()
		if err != nil {
			return fmt.Errorf("backup :%w", err)
		}

		log.Info().Msgf("success backup to %s", outputPath)
		return nil
	}

	wrk := worker.New(bak, worker.WithCronTab(cfg.CronTab))
	if err = wrk.Run(); err != nil {
		return err
	}

	return nil
}
