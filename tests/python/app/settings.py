from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    QDRANT_URL: str = 'http://qdrant:6333'
    OLLAMA_URL: str = 'http://localhost:11434'
    OLLAMA_EMBED_MODEL: str = 'embeddinggemma'
    OLLAMA_LLM_MODEL: str = 'gemma2:2b'
    QDRANT_VECTOR_SIZE: int = 512

    model_config = SettingsConfigDict(env_file=".env")

settings = Settings()
