from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from app.qdrant_client import get_client, ensure_collection
from app.settings import settings
import requests

router = APIRouter(prefix='/rag', tags=['rag'])

class QueryIn(BaseModel):
    query: str
    top_k: int = 4

@router.post('/query')
def query_rag(payload: QueryIn):
    client = get_client()
    ensure_collection('trusts')
    try:
        res = client.search(collection_name='trusts', query_vector=[0.0]*settings.QDRANT_VECTOR_SIZE, limit=payload.top_k)
        docs = [p.payload for p in res]
    except Exception:
        docs = []
    # call Ollama local LLM (user runs separately)
    try:
        url = f"{settings.OLLAMA_URL}/completions?model={settings.OLLAMA_LLM_MODEL}"
        resp = requests.post(url, json={"prompt": payload.query, "context": docs}, timeout=10)
        resp.raise_for_status()
        return resp.json()
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Ollama call failed: {e}")
