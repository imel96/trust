import requests
from app.settings import settings

def embed_text(text: str):
    try:
        url = f"{settings.OLLAMA_URL}/embed?model={settings.OLLAMA_EMBED_MODEL}"
        resp = requests.post(url, json={"input": text}, timeout=10)
        resp.raise_for_status()
        data = resp.json()
        if isinstance(data, dict) and 'embedding' in data:
            return data['embedding']
        if isinstance(data, list) and len(data)>0 and 'embedding' in data[0]:
            return data[0]['embedding']
        # fallback
        return [0.0]*settings.QDRANT_VECTOR_SIZE
    except Exception:
        return [0.0]*settings.QDRANT_VECTOR_SIZE
