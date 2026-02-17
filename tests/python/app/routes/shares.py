from datetime import datetime
from decimal import Decimal
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
def buy_shares(code: str, qty: int, price: float=10.0):
    commission = 1.0
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    total_price = qty * price
    payload['cash'] = payload.get('cash',0.0) - total_price - commission
    payload['cost'] = payload.get('cost',0.0) - commission
    payload['shares'] = total_price
    portfolio = payload.get('portfolio',{})
    portfolio[code] = portfolio.get(code, 0) + qty
    payload['transactions'] = create_transactions_payload(tp.payload, code, qty, price, total_price)
    payload['portfolio'] = portfolio
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='bought')

@router.post('/sell', response_model=SimpleResponse)
def sell_shares(code: str, qty: int, price: float=10.0):
    commission = 1.0
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    portfolio = payload.get('portfolio',{})
    current = portfolio.get(code,0)
    if current < qty:
        raise HTTPException(status_code=400, detail='Not enough shares')
    portfolio[code] = current - qty
    total_price = qty * price
    payload['portfolio'] = portfolio
    payload['cost'] = payload.get('cost', 0.0) - commission
    payload['cash'] = payload.get('cash', 0.0) + total_price - commission
    payload['shares'] = payload.get('shares',0.0) - total_price
    payload['transactions'] = create_transactions_payload(tp.payload, code, -qty, price, total_price)
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='sold')

def create_transactions_payload(payload: dict, security: str, units: int, price_per_unit: float, amount: float) -> dict:
    transactions = payload.get('transactions', [])
    transactions.append({
        'transaction_date': datetime.now(),
        'security': security,
        'units': units,
        'price_per_unit': Decimal(price_per_unit),
        'amount': Decimal(amount)
    })
    return transactions
