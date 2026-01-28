from fastapi import APIRouter, HTTPException
from app.qdrant_client import get_client, ensure_collection
from app.rag_utils import embed_text
from app.entities import SimpleResponse

router = APIRouter(prefix='/income', tags=['income'])

def _get_trust_point(client):
    ensure_collection('trusts')
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    if pts:
        return pts[0]
    return None

@router.post('/dividend', response_model=SimpleResponse)
def pay_dividend(dollar_per_share: float):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    portfolio = payload.get('portfolio', {})
    total_shares = sum(portfolio.values()) if portfolio else 0
    amount = total_shares * dollar_per_share
    payload['income'] = payload.get('income', 0.0) + amount
    payload['cash'] = payload.get('cash', 0.0) + amount
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail=str(amount))
