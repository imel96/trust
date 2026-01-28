from qdrant_client import QdrantClient
from app.settings import settings

_client = None

def get_client():
    global _client
    if _client is None:
        url = settings.QDRANT_URL.replace('http://','').replace('https://','').split(':')[0]
        _client = QdrantClient(url=settings.QDRANT_URL, prefer_grpc=False)
    return _client

def ensure_collection(name, vector_size=None):
    client = get_client()
    try:
        client.get_collection(collection_name=name)
    except Exception:
        client.create_collection(collection_name=name, vectors_config={'size': settings.QDRANT_VECTOR_SIZE, 'distance': 'Cosine'})
