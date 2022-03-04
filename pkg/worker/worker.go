package worker

import (
	"time"

	"github.com/fahmifan/dbback/pkg/model"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

type Worker struct {
	cronScheduler *gocron.Scheduler
	backuper      model.Backupper
}

func New(backuper model.Backupper) *Worker {
	return &Worker{
		backuper:      backuper,
		cronScheduler: gocron.NewScheduler(time.Local),
	}
}

func (w *Worker) Run() error {
	w.registerPeriodicBackup()

	log.Info().Msg("start cron job")
	w.cronScheduler.StartBlocking()
	return nil
}

func (w *Worker) registerPeriodicBackup() error {
	_, err := w.cronScheduler.Every(1).Day().At("00:00").Do(func() {
		outputPath, err := w.backuper.Backup()
		if err != nil {
			log.Error().Err(err).Msg("periodic")
			return
		}
		log.Info().Msgf("success backup to %s", outputPath)
	})
	if err != nil {
		return err
	}
	return nil
}
