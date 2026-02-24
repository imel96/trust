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
    payload['transactions'] = create_transactions_payload(payload['transactions'], code, qty, price, -total_price, 'cash', 'shares')
    payload['transactions'] = create_transactions_payload(payload['transactions'], None, 0, 0, -commission, 'cash', 'cost')
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
    # XXX
    portfolio = payload.get('portfolio',{})
    current = portfolio.get(code,0)
    if current['units'] < qty:
        raise HTTPException(status_code=400, detail='Not enough shares')
    # /XXX
    total_price = qty * price
    payload['transactions'] = create_transactions_payload(payload['transactions'], code, -qty, price, total_price, 'cash', 'shares')
    payload['transactions'] = create_transactions_payload(payload['transactions'], None, 0, 0, -commission, 'cash', 'cost')
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='sold')

def create_transactions_payload(transactions: list, security: str, units: int, price_per_unit: float, amount: float, cash_account: str, offset_account: str) -> list:
    transactions.append({
        'transaction_date': datetime.now(),
        'security': security,
        'units': units,
        'price_per_unit': Decimal(price_per_unit),
        'amount': Decimal(amount),
        'cash_account': cash_account,
        'offset_account': offset_account,
    })
    return transactions
