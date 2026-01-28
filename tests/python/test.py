from qdrant_client import QdrantClient
client = QdrantClient(host="qdrant", port=6333)
pts, _ = client.scroll(collection_name="trusts", limit=10)
print("Trusts found:", pts)

