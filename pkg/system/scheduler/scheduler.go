package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/hidevopsio/hiboot/pkg/log"
)

type Scheduler struct {
	*gocron.Scheduler
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Scheduler: gocron.NewScheduler(time.UTC),
	}
}

func (s *Scheduler) RunOnce(task func())  {
	job, _ := s.Every(100).Milliseconds().Do(task)
	job.LimitRunsTo(1)
	s.StartAsync()
}

func (s *Scheduler) Run(tag *string, limit *int, every *int, unit *string, atTime *string, delay *int64, task func())  {
	defaultEvery := 1
	if every == nil {
		every = &defaultEvery
	}
	_ = s.Every(*every)

	if unit != nil {
		schedulingUnit := gocron.Milliseconds
		switch *unit {
		case "milliseconds":
			schedulingUnit = gocron.Milliseconds
		case "seconds":
			schedulingUnit = gocron.Seconds
		case "minutes":
			schedulingUnit = gocron.Minutes
		case "hours":
			schedulingUnit = gocron.Hours
		case "days":
			schedulingUnit = gocron.Days
		case "weeks":
			schedulingUnit = gocron.Weeks
		case "months":
			schedulingUnit = gocron.Months
		default:
			schedulingUnit = gocron.Seconds
			log.Warn("invalid unit for scheduler, use seconds as default")
		}
		s.SetUnit(schedulingUnit)
	}
	if atTime != nil {
		s.At(*atTime)
	}
	if limit != nil {
		s.LimitRunsTo(*limit)
	}
	if tag != nil {
		s.Tag(*tag)
	}
	if delay == nil {
		s.RunAll()
	} else {
		s.RunAllWithDelay(time.Duration(*delay))
	}

	_, err := s.Do(task)
	if err != nil {
		log.Error(err)
		return
	}

	s.StartAsync()
}

func (s *Scheduler) RunWithExpr(tag *string, expressions *string, task func())  {

	s.Cron(*expressions)
	if tag != nil {
		s.Tag(*tag)
	}
	_, err := s.Do(task)
	if err == nil {
		s.StartAsync()
	}
}
