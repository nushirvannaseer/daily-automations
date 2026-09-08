package scheduler

import (
	"log"
	"github.com/robfig/cron/v3"
)

// Scheduler wraps the cron library.
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new scheduler instance.
func NewScheduler() *Scheduler {
	// Using the standard cron format (minute, hour, dom, month, dow)
	return &Scheduler{
		cron: cron.New(),
	}
}

// AddJob adds a job to the scheduler.
// schedule is a standard cron expression (e.g., "0 0 * * *" for daily).
func (s *Scheduler) AddJob(schedule string, name string, job func()) error {
	_, err := s.cron.AddFunc(schedule, func() {
		log.Printf("Starting scheduled job: %s\n", name)
		job()
		log.Printf("Finished scheduled job: %s\n", name)
	})
	
	if err != nil {
		return err
	}
	
	log.Printf("Registered job '%s' with schedule '%s'\n", name, schedule)
	return nil
}

// Start runs the scheduler. This method is non-blocking.
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started...")
}

// Stop stops the scheduler.
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped.")
}
