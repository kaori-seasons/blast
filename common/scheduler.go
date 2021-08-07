package common

import (
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

type scheduler struct {
	routine     *cron.Cron
	builtInJobs map[string]cron.EntryID // name => entry
}

type JobWrapper struct {
	Name   string
	Handle func() error
}

func (jw JobWrapper) Run() {
	err := jw.Handle()
	if err != nil {
		log.Warnf("run job(%s) failed. err = %s", jw.Name, err)
	}
}

func AddRoutineJob(spec string, job JobWrapper) error {
	if schedulerInstance == nil {
		initScheduler()
	}

	id, e := schedulerInstance.routine.AddJob(spec, job)
	if e != nil {
		return e
	}
	log.Infof("AddRoutineJob name = %s spec = %s", job.Name, spec)
	schedulerInstance.builtInJobs[job.Name] = id
	return nil
}

func StartSchedule() {
	if schedulerInstance != nil {
		schedulerInstance.routine.Start()
	}
}

var schedulerInstance *scheduler

func initScheduler() {
	schedulerInstance = &scheduler{
		routine:     cron.New(cron.WithSeconds(), cron.WithLogger(cron.PrintfLogger(newCronLogger()))),
		builtInJobs: make(map[string]cron.EntryID),
	}
}

type cronLog struct{}

func newCronLogger() *cronLog {
	return &cronLog{}
}

func (l *cronLog) Printf(msg string, v ...interface{}) {
	log.Infof(msg, v...)
}
