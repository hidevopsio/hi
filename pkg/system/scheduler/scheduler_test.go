package scheduler_test

import (
	"sync"
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/system/scheduler"
)

var waitGroup = sync.WaitGroup{}

type myService struct {
	at.Scheduled `limit:"1" every:"1" unit:"seconds"`
}

func (s *myService) Run()  {
	log.Info("Running Task")
	waitGroup.Done()
}

type myLimitService struct {
	at.Scheduled `limit:"2"`
}

func (s *myLimitService) Run()  {
	log.Info("Running Task")
	waitGroup.Done()
}

func TestSchedulerService(t *testing.T) {
	svc := new(myService)
	waitGroup.Add(1)

	s := scheduler.NewScheduler()
	s.RunOnce(svc.Run)

	waitGroup.Wait()
	log.Info("Done")
}

func TestRunOnce(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	s := scheduler.NewScheduler()
	s.RunOnce(func() {
		wg.Done()
		log.Info("Running Task")
	})

	wg.Wait()
	log.Info("Done")
}


func TestExampleJob_LimitRunsTo(t *testing.T) {
	count := 2
	wg := sync.WaitGroup{}
	wg.Add(count)
	s := gocron.NewScheduler(time.UTC)
	job, _ := s.Every(1).Second().Do(func() {
		wg.Done()
		log.Info("Running Task")
	})
	job.LimitRunsTo(count)
	s.StartAsync()
	wg.Wait()
}