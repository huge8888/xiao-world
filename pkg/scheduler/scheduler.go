package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/processor"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/publishers"
	"github.com/xpzouying/xiaohongshu-mcp/pkg/types"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

// Scheduler manages scheduled posts
type Scheduler struct {
	processor  *processor.Processor
	publishers map[types.Platform]publishers.Publisher
	jobs       map[string]*types.ScheduledJob
	mu         sync.RWMutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewScheduler creates a new scheduler
func NewScheduler(proc *processor.Processor, pubs map[types.Platform]publishers.Publisher) *Scheduler {
	return &Scheduler{
		processor:  proc,
		publishers: pubs,
		jobs:       make(map[string]*types.ScheduledJob),
		stopCh:     make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go s.run()
	logrus.Info("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
	logrus.Info("Scheduler stopped")
}

// run is the main scheduler loop
func (s *Scheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.processScheduledJobs()
		}
	}
}

// processScheduledJobs processes all scheduled jobs that are due
func (s *Scheduler) processScheduledJobs() {
	s.mu.RLock()
	dueJobs := make([]*types.ScheduledJob, 0)
	now := time.Now()

	for _, job := range s.jobs {
		if job.Status == types.JobStatusPending && job.ScheduledAt.Before(now) {
			dueJobs = append(dueJobs, job)
		}
	}
	s.mu.RUnlock()

	// Process due jobs
	for _, job := range dueJobs {
		s.executeJob(job)
	}
}

// executeJob executes a scheduled job
func (s *Scheduler) executeJob(job *types.ScheduledJob) {
	s.mu.Lock()
	job.Status = types.JobStatusRunning
	s.mu.Unlock()

	logrus.Infof("Executing scheduled job: %s", job.ID)

	// Get feed detail from Xiaohongshu
	// Note: This requires Xiaohongshu service to be available
	// For now, we'll assume feed details are already stored in the job
	// In a production system, you'd fetch the feed here

	var results []types.PublishResult

	// Publish to each platform
	for _, platform := range job.Platforms {
		publisher, exists := s.publishers[platform]
		if !exists || !publisher.IsEnabled() {
			logrus.Warnf("Publisher for platform %s not available or disabled", platform)
			results = append(results, types.PublishResult{
				Platform:  platform,
				Success:   false,
				Error:     fmt.Sprintf("publisher not available or disabled"),
				Timestamp: time.Now(),
			})
			continue
		}

		// Note: In production, you'd fetch the feed detail and process it
		// For now, this is a placeholder
		logrus.Infof("Publishing to %s", platform)

		// Process content for platform
		// content, err := s.processor.Process(feedDetail, platform)
		// if err != nil {
		// 	logrus.Errorf("Failed to process content for %s: %v", platform, err)
		// 	continue
		// }

		// Publish content
		// result, err := publisher.Publish(content)
		// if err != nil {
		// 	logrus.Errorf("Failed to publish to %s: %v", platform, err)
		// }
		// results = append(results, *result)
	}

	s.mu.Lock()
	job.Status = types.JobStatusCompleted
	job.Results = results
	s.mu.Unlock()

	logrus.Infof("Completed scheduled job: %s", job.ID)
}

// ScheduleJob schedules a new job
func (s *Scheduler) ScheduleJob(feedID string, platforms []types.Platform, scheduledAt time.Time) (string, error) {
	if scheduledAt.Before(time.Now()) {
		return "", fmt.Errorf("scheduled time is in the past")
	}

	if len(platforms) == 0 {
		return "", fmt.Errorf("no platforms specified")
	}

	job := &types.ScheduledJob{
		ID:          uuid.New().String(),
		FeedID:      feedID,
		Platforms:   platforms,
		ScheduledAt: scheduledAt,
		Status:      types.JobStatusPending,
		Results:     make([]types.PublishResult, 0),
	}

	s.mu.Lock()
	s.jobs[job.ID] = job
	s.mu.Unlock()

	logrus.Infof("Scheduled job %s for feed %s at %s", job.ID, feedID, scheduledAt)

	return job.ID, nil
}

// GetJob retrieves a job by ID
func (s *Scheduler) GetJob(jobID string) (*types.ScheduledJob, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return nil, fmt.Errorf("job not found: %s", jobID)
	}

	return job, nil
}

// ListJobs lists all jobs
func (s *Scheduler) ListJobs() []*types.ScheduledJob {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make([]*types.ScheduledJob, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

// CancelJob cancels a scheduled job
func (s *Scheduler) CancelJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	if job.Status == types.JobStatusRunning {
		return fmt.Errorf("cannot cancel running job")
	}

	if job.Status == types.JobStatusCompleted {
		return fmt.Errorf("cannot cancel completed job")
	}

	job.Status = types.JobStatusCancelled

	logrus.Infof("Cancelled job: %s", jobID)

	return nil
}

// DeleteJob deletes a job
func (s *Scheduler) DeleteJob(jobID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, exists := s.jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	if job.Status == types.JobStatusRunning {
		return fmt.Errorf("cannot delete running job")
	}

	delete(s.jobs, jobID)

	logrus.Infof("Deleted job: %s", jobID)

	return nil
}

// PublishNow publishes content to specified platforms immediately
func (s *Scheduler) PublishNow(feed *xiaohongshu.FeedDetail, platforms []types.Platform) ([]types.PublishResult, error) {
	if len(platforms) == 0 {
		return nil, fmt.Errorf("no platforms specified")
	}

	var results []types.PublishResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Publish to each platform concurrently
	for _, platform := range platforms {
		wg.Add(1)
		go func(p types.Platform) {
			defer wg.Done()

			publisher, exists := s.publishers[p]
			if !exists || !publisher.IsEnabled() {
				logrus.Warnf("Publisher for platform %s not available or disabled", p)
				mu.Lock()
				results = append(results, types.PublishResult{
					Platform:  p,
					Success:   false,
					Error:     fmt.Sprintf("publisher not available or disabled"),
					Timestamp: time.Now(),
				})
				mu.Unlock()
				return
			}

			// Process content for platform
			content, err := s.processor.Process(feed, p)
			if err != nil {
				logrus.Errorf("Failed to process content for %s: %v", p, err)
				mu.Lock()
				results = append(results, types.PublishResult{
					Platform:  p,
					Success:   false,
					Error:     fmt.Sprintf("failed to process content: %v", err),
					Timestamp: time.Now(),
				})
				mu.Unlock()
				return
			}

			// Publish content
			result, err := publisher.Publish(content)
			if err != nil {
				logrus.Errorf("Failed to publish to %s: %v", p, err)
			}

			mu.Lock()
			results = append(results, *result)
			mu.Unlock()
		}(platform)
	}

	wg.Wait()

	return results, nil
}
