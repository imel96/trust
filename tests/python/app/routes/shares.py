from fastapi import APIRouter, HTTPException
from app.qdrant_client import get_client, ensure_collection
from app.rag_utils import embed_text
from app.entities import SimpleResponse

router = APIRouter(prefix='/shares', tags=['shares'])

def _get_trust_point(client):
    ensure_collection('trusts')
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    if pts:
        return pts[0]
    return None

@router.post('/buy', response_model=SimpleResponse)
def buy_shares(symbol: str, qty: int, price: float=0.1):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    total_price = qty * price
    payload['cash'] = payload.get('cash',0.0) - total_price
    payload['cost'] = payload.get('cost',0.0) + 1.0
    payload['shares'] = total_price
    portfolio = payload.get('portfolio',{})
    portfolio[symbol] = portfolio.get(symbol,0) + qty
    payload['portfolio'] = portfolio
    print("shares.py", payload)
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='bought')

@router.post('/sell', response_model=SimpleResponse)
def sell_shares(symbol: str, qty: int, price: float=10.0):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    portfolio = payload.get('portfolio',{})
    current = portfolio.get(symbol,0)
    if current < qty:
        raise HTTPException(status_code=400, detail='Not enough shares')
    portfolio[symbol] = current - qty
    total_price = qty * price
    payload['portfolio'] = portfolio
    payload['cost'] = payload.get('cost',0.0) + 1.0
    payload['cash'] = payload.get('cash',0.0) + total_price - payload['cost']
    payload['shares'] = payload.get('shares',0.0) - total_price
    print("shares.py", payload)
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='sold')
