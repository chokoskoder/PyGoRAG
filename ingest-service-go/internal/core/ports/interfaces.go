package ports

import (
	"context"
	"github.com/chokoskoder/PyGoRAG/internal/core/domain"
)

type DocumentLoader interface {
	Load(ctx context.Context, job domain.IngestJob) ([]domain.Document, error)
}

type VectorRepository interface {
	Store(ctx context.Context, docs []domain.Document) error
}

type Embedder interface {
	EmbedQuery(ctx context.Context, text string) ([]float32, error)
	// Usually handled by the vector repo adapter, but exposed if needed
}