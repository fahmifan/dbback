package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/fahmifan/dbback/pkg/backuper"
	"github.com/fahmifan/dbback/pkg/config"
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

func run(args []string) error {
	cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
	var (
		outputPath string
		dbName     string
		dbDriver   string
	)

	cmd.StringVar(&dbName, "dbname", "", `--dbname foobar`)
	cmd.StringVar(&dbDriver, "driver", "", `--driver [mysql, c, postgres]`)

	if err := cmd.Parse(args[1:]); err != nil {
		return fmt.Errorf("parse args: %w", err)
	}

	cfg, err := config.Load("./config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	switch dbDriver {
	default:
		return errors.New("invalid driver, should be [mysql, postgres]")
	case "postgres":
		bak := backuper.NewPostgre(&backuper.PostgreCfg{
			OutDir:   cfg.OutDir,
			User:     cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			DBName:   dbName,
		})
		outputPath, err = bak.Backup()
	case "mysql":
		bak := backuper.NewMySQL(&backuper.MySQLCfg{
			OutDir:   cfg.OutDir,
			User:     cfg.MySQL.User,
			Password: cfg.MySQL.Password,
			Host:     cfg.MySQL.Host,
			Port:     cfg.MySQL.Port,
			DBName:   dbName,
		})
		outputPath, err = bak.Backup()
	}
	if err != nil {
		return fmt.Errorf("backup :%w", err)
	}

	log.Info().Msgf("success backup to %s", outputPath)
	return nil
}
