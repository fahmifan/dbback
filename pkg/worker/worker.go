package worker

import (
	"time"

	"github.com/fahmifan/dbback/pkg/model"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// at every 00:00
const defaultCronTab = "0 0 * * *"

type Worker struct {
	cronScheduler *gocron.Scheduler
	backuper      model.Backupper
	cronTab       string
}

type opt func(w *Worker)

func WithCronTab(cronTab string) opt {
	return func(w *Worker) {
		if cronTab == "" {
			cronTab = defaultCronTab
		}
		w.cronTab = cronTab
	}
}

func New(backuper model.Backupper, opts ...opt) *Worker {
	wrk := &Worker{
		backuper:      backuper,
		cronScheduler: gocron.NewScheduler(time.Local),
		cronTab:       defaultCronTab,
	}

	for _, opt := range opts {
		opt(wrk)
	}

	return wrk
}

func (w *Worker) Run() error {
	w.registerPeriodicBackup()

	log.Info().Msg("start cron job")
	w.cronScheduler.StartBlocking()
	return nil
}

func (w *Worker) registerPeriodicBackup() error {
	_, err := w.cronScheduler.Cron(w.cronTab).Do(func() {
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
