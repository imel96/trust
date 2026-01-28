from fastapi import APIRouter, HTTPException
from app.qdrant_client import get_client, ensure_collection
from app.rag_utils import embed_text
from app.entities import SimpleResponse
from qdrant_client import QdrantClient

router = APIRouter(prefix='/trust-unit', tags=['trust_unit'])

def _get_trust_point(client):
    ensure_collection('trusts')
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    if pts:
        return pts[0]
    return None

@router.post('/invest', response_model=SimpleResponse)
def invest(person_id: str, amount: float):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    print(tp)
    payload = tp.payload
    payload['cash'] = payload.get('cash',0.0) + amount
    payload['units_outstanding'] = payload.get('units_outstanding',0) + 1
    payload['nav_per_unit'] = payload['cash']/payload['units_outstanding'] if payload['units_outstanding']>0 else 0.0
    # update person units
    ensure_collection('persons')
    pts, x = client.scroll(collection_name='persons', limit=10)
    print(pts)
    print(x)
    for p in pts:
        if p.payload.get('name') == person_id:
            p_payload = p.payload
            p_payload['units'] = p_payload.get('units',0)+1
            client.upsert(collection_name='persons', points=[{'id':p.id, 'vector': embed_text(p_payload.get('text','')), 'payload': p_payload}])
            break
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='invested')

@router.post('/redeem', response_model=SimpleResponse)
def redeem(person_id: str, amount: float):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    payload['cash'] = max(0.0, payload.get('cash',0.0) - amount)
    payload['units_outstanding'] = max(0, payload.get('units_outstanding',0)-1)
    payload['nav_per_unit'] = payload['cash']/payload['units_outstanding'] if payload['units_outstanding']>0 else 0.0
    # update person units
    ensure_collection('persons')
    pts, _ = client.scroll(collection_name='persons', limit=10)
    for p in pts:
        if p.payload.get('name') == person_id:
            p_payload = p.payload
            p_payload['units'] = max(0, p_payload.get('units',0)-1)
            client.upsert(collection_name='persons', points=[{'id':p.id, 'vector': embed_text(p_payload.get('text','')), 'payload': p_payload}])
            break
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='redeemed')
