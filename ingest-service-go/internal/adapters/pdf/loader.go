package pdf

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"github.com/ledongthuc/pdf"
	"github.com/tmc/langchaingo/textsplitter"

	"github.com/chokoskoder/PyGoRAG/internal/core/domain"
)

type WebPDFLoader struct{}

func NewLoader() *WebPDFLoader{
	return &WebPDFLoader{}
}
func (l *WebPDFLoader) Load(ctx context.Context , job domain.IngestJob) ([]domain.Document, error){
	//Fetch 
	req,_ := http.NewRequestWithContext(ctx , "GET" , job.SourceURL , nil)
	resp, err := http.DefaultClient.Do(req)
	if err !=nil {
		return nil , fmt.Errorf("download failed : %w" , err)
	}

	body , err := io.ReadAll(resp.Body)
	if err != nil {return nil , err}

	//Parse
	reader := bytes.NewReader(body)
	r, err := pdf.NewReader(reader , int64(len(body)))

	if err!= nil  {return nil , err}

	var rawText bytes.Buffer
	for i := 1 ; i<= r.NumPage(); i++{
		p := r.Page(i)
		if !p.V.IsNull(){
			t,_ := p.GetPlainText(nil)
			rawText.WriteString(t)
		}
	}
	//Chunk
	splitter := textsplitter.NewRecursiveCharacter()
	splitter.ChunkSize = 800
	splitter.ChunkOverlap = 100
	chunks, err := splitter.SplitText(rawText.String())
	if err != nil { return nil, err }

	// 4. Map to Domain
	var docs []domain.Document
	for _, c := range chunks {
		docs = append(docs, domain.Document{
			Content: c,
			Source:  job.SourceURL,
			Metadata: map[string]interface{}{
				"title": job.Title,
				"url":   job.SourceURL,
			},
		})
	}
	return docs, nil
}
