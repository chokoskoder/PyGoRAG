package main

import (
	"context"
	"github.com/chokoskoder/PyGoRAG/internal/adapters/pdf"
	"github.com/chokoskoder/PyGoRAG/internal/adapters/qdrant"
	"github.com/chokoskoder/PyGoRAG/internal/core/domain"
	"github.com/chokoskoder/PyGoRAG/internal/core/services"
	"log/slog"
	"os"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms/ollama"
)

func main() {
	// 1. Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	
	// 2. Init Adapters (Infra)
	llm, err := ollama.New(ollama.WithModel("nomic-embed-text"))
	if err != nil { logger.Error("Failed to init Ollama", "err", err); os.Exit(1) }
	
	embedder, err := embeddings.NewEmbedder(llm)
	if err != nil { logger.Error("Failed to init Embedder", "err", err); os.Exit(1) }

	repo, err := qdrant.NewRepository("http://localhost:6333", "production_data", embedder, logger)
	if err != nil { logger.Error("Failed to init Qdrant", "err", err); os.Exit(1) }

	loader := pdf.NewLoader()

	// 3. Init Service (Business Logic)
	svc := services.NewIngestionService(loader, repo, logger, 3)

	// 4. Run
	jobs := []domain.IngestJob{
		{Title: "Attention Is All You Need", SourceURL: "https://arxiv.org/pdf/1706.03762.pdf"},
		{Title: "Llama 2 Paper", SourceURL: "https://arxiv.org/pdf/2307.09288.pdf"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := svc.Run(ctx, jobs); err != nil {
		logger.Error("Ingestion finished with errors")
	} else {
		logger.Info("Ingestion pipeline completed successfully")
	}
}