package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"sales_analytics/config"
	"sales_analytics/pkg/repository"

	"github.com/robfig/cron/v3"
)

// Scheduler manages cron jobs for data refresh
type Scheduler struct {
	cron      *cron.Cron
	jobID     cron.EntryID
	jobLock   sync.Mutex
	repo      *repository.MongoRepository
	config    *config.Config
	isRunning bool
}

// NewScheduler creates a new scheduler instance
func NewScheduler(repo *repository.MongoRepository, cfg *config.Config) *Scheduler {
	return &Scheduler{
		cron:      cron.New(),
		repo:      repo,
		config:    cfg,
		isRunning: false,
	}
}

// Start starts the cron scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Cron scheduler started")
}

// Stop stops the cron scheduler and removes all jobs
func (s *Scheduler) Stop() {
	s.jobLock.Lock()
	defer s.jobLock.Unlock()

	if s.jobID != 0 {
		s.cron.Remove(s.jobID)
		s.jobID = 0
	}

	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Cron scheduler stopped and all jobs cleaned up")
}

// CreateJob creates or replaces the data refresh cron job
func (s *Scheduler) CreateJob(interval string) error {
	s.jobLock.Lock()
	defer s.jobLock.Unlock()

	// Remove existing job if exists
	if s.jobID != 0 {
		s.cron.Remove(s.jobID)
		log.Printf("Removed existing cron job (ID: %d)", s.jobID)
		s.jobID = 0
	}

	// Validate interval format
	cronExpr := "@every " + interval

	// Add the new cron job
	id, err := s.cron.AddFunc(cronExpr, s.executeDataRefresh)
	if err != nil {
		return fmt.Errorf("failed to create cron job: %w", err)
	}

	s.jobID = id
	log.Printf("Created new cron job (ID: %d) with interval: %s", id, interval)

	return nil
}

// DeleteJob removes the current cron job
func (s *Scheduler) DeleteJob() error {
	s.jobLock.Lock()
	defer s.jobLock.Unlock()

	if s.jobID == 0 {
		return fmt.Errorf("no active cron job to delete")
	}

	s.cron.Remove(s.jobID)
	log.Printf("Deleted cron job (ID: %d)", s.jobID)
	s.jobID = 0

	return nil
}

// GetJobStatus returns the current job status
func (s *Scheduler) GetJobStatus() map[string]interface{} {
	s.jobLock.Lock()
	defer s.jobLock.Unlock()

	status := map[string]interface{}{
		"active": s.jobID != 0,
		"job_id": s.jobID,
	}

	if s.jobID != 0 {
		entry := s.cron.Entry(s.jobID)
		if entry.ID != 0 {
			status["next_run"] = entry.Next
			status["previous_run"] = entry.Prev
		}
	}

	return status
}

// executeDataRefresh performs the actual data refresh
func (s *Scheduler) executeDataRefresh() {
	// Prevent concurrent executions
	s.jobLock.Lock()
	if s.isRunning {
		log.Println("Data refresh already running, skipping this execution")
		s.jobLock.Unlock()
		return
	}
	s.isRunning = true
	s.jobLock.Unlock()

	defer func() {
		s.jobLock.Lock()
		s.isRunning = false
		s.jobLock.Unlock()
	}()

	log.Println("Cron job triggered: Starting data refresh...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	loader := repository.NewDataLoader(s.repo, s.config.WorkerPoolSize)

	if err := loader.LoadCSV(ctx, s.config.CSVFilePath); err != nil {
		log.Printf("Cron job failed: %v", err)
	} else {
		log.Println("Cron job completed successfully")
	}
}
