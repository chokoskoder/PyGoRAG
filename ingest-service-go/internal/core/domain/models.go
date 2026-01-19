package domain

type Document struct {
	ID       string
	Content  string
	Source   string
	Metadata map[string]interface{}
}

type IngestJob struct {
	SourceURL string
	Title     string
}