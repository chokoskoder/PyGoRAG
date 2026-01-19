from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    QDRANT_URL: str = "http://localhost:6333"
    COLLECTION_NAME: str = "production_data"
    OLLAMA_BASE_URL: str = "http://localhost:11434"
    OLLAMA_MODEL: str = "mistral"
    EMBEDDING_MODEL: str = "nomic-embed-text:latest"

    class Config:
        env_file = ".env"

settings = Settings()