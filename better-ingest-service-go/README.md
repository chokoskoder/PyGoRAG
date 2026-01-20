# Ingest service written in go 
This will be a re writing of the previous go service which was written by AI to ensure better decoupling and introduce a few new features like :

## Exponential Backoff with Jitter
This will ensure that when more than 5 PDFs are loaded together our ELT pipeline doesnt get overwhelmed and stops working completely. It will introduce a **Loading** time before the PDF enters the ELT pipeline.
The backoff is planned to be 15s with a jitter of ±25% to ensure that no two pdfs enter at the same time which create an edge case of either one of the PDFs having to wait almost 30s

## Overlapping/Sliding Window Chunking
This ingestion service will introduce a better method of context preservation by replacing the already fixed chunking process by **Overlapping** chunking. 

## Better decoupling and writing industry ready Production grade go
This code will include:
- Logging
    - We will introduce the logging of chunk creation : {chunk index + filename}
- Config Management 
- Production ENV and Developer ENV
- will be written following idiomatic go practices

