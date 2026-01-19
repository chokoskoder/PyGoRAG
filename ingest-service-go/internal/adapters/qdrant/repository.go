package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/chokoskoder/PyGoRAG/internal/core/domain"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/vectorstores/qdrant"
)

type Repository struct {
	store  *qdrant.Store // We keep this as a pointer for efficiency
	logger *slog.Logger
}

// NewRepository creates a qdrant instance AND ensures the collection exists
func NewRepository(baseURL, collection string, embedder embeddings.Embedder, logger *slog.Logger) (*Repository, error) {
	// 1. Check Schema (using string URL for HTTP check)
	if err := ensureCollectionExists(baseURL, collection, 768, logger); err != nil {
		return nil, fmt.Errorf("failed to migrate database schema: %w", err)
	}

	// 2. Parse URL (Fixes the "url type" error)
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid qdrant url: %w", err)
	}

	// 3. Create Store
	s, err := qdrant.New(
		qdrant.WithURL(*parsedURL), // Fix: Pass the dereferenced URL struct
		qdrant.WithCollectionName(collection),
		qdrant.WithEmbedder(embedder),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to qdrant: %w", err)
	}

	// 4. Return Repository (Fixes the "structure literal" error)
	// We take the address of 's' (&s) because our struct expects a pointer
	return &Repository{store: &s, logger: logger}, nil
}

func (r *Repository) Store(ctx context.Context, docs []domain.Document) error {
	r.logger.Info("Storing batch in Vector DB", "count", len(docs))

	lcDocs := make([]schema.Document, len(docs))
	for i, d := range docs {
		lcDocs[i] = schema.Document{
			PageContent: d.Content,
			Metadata:    d.Metadata,
		}
	}

	_, err := r.store.AddDocuments(ctx, lcDocs)
	return err
}

// ensureCollectionExists checks if the collection exists, and creates it if not.
func ensureCollectionExists(baseURL, name string, vectorSize int, logger *slog.Logger) error {
	checkURL := fmt.Sprintf("%s/collections/%s", baseURL, name)
	resp, err := http.Get(checkURL)
	if err != nil {
		return fmt.Errorf("failed to check collection status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		logger.Info("Collection exists", "name", name)
		return nil
	}

	logger.Info("Collection missing. Performing Auto-Migration...", "name", name)
	
	createURL := fmt.Sprintf("%s/collections/%s", baseURL, name)
	payload := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}
	
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PUT", createURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	createResp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	defer createResp.Body.Close()

	if createResp.StatusCode != 200 {
		return fmt.Errorf("failed to create collection, status: %d", createResp.StatusCode)
	}

	logger.Info("Schema Migration Complete. Collection created.", "name", name)
	return nil
}