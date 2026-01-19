package services

import (
	"context"
	"github.com/chokoskoder/PyGoRAG/internal/core/domain"
	"github.com/chokoskoder/PyGoRAG/internal/core/ports"
	"log/slog"
	"sync"
)

type IngestionService struct {
	loader   ports.DocumentLoader
	repo     ports.VectorRepository
	logger   *slog.Logger
	workers  int
}

func NewIngestionService(l ports.DocumentLoader, r ports.VectorRepository, logger *slog.Logger, workers int) *IngestionService {
	return &IngestionService{loader: l, repo: r, logger: logger, workers: workers}
}

func (s *IngestionService) Run(ctx context.Context, jobs []domain.IngestJob) error {
	jobChan := make(chan domain.IngestJob, len(jobs))
	errChan := make(chan error, len(jobs))
	var wg sync.WaitGroup

	s.logger.Info("Starting Ingestion Worker Pool", "workers", s.workers, "jobs", len(jobs))

	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobChan {
				s.logger.Debug("Worker started job", "worker_id", id, "job", job.Title)
				
				docs, err := s.loader.Load(ctx, job)
				if err != nil {
					s.logger.Error("Failed to load document", "worker_id", id, "error", err)
					errChan <- err
					continue
				}

				if err := s.repo.Store(ctx, docs); err != nil {
					s.logger.Error("Failed to store documents", "worker_id", id, "error", err)
					errChan <- err
					continue
				}
				s.logger.Info("Worker completed job", "worker_id", id, "job", job.Title)
			}
		}(i)
	}

	for _, j := range jobs {
		jobChan <- j
	}
	close(jobChan)
	wg.Wait()
	close(errChan)

	// In a real system, you might return aggregated errors or a percentage success rate
	return nil
}