from fastapi import APIRouter, HTTPException
from app.qdrant_client import get_client, ensure_collection
from app.rag_utils import embed_text
from app.entities import SimpleResponse

router = APIRouter(prefix='/trusts', tags=['trusts'])

def _get_trust_point(client):
    ensure_collection('trusts')
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    if pts:
        return pts[0]
    return None

@router.post('/add_member', response_model=SimpleResponse)
def add_member(person_id: str):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    members = payload.get('members', [])
    members.append(person_id)
    payload['members'] = members
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='member added')

@router.post('/add_broker', response_model=SimpleResponse)
def add_broker(broker_id: str):
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        raise HTTPException(status_code=400, detail='No trust found')
    payload = tp.payload
    brokers = payload.get('brokers', [])
    brokers.append(broker_id)
    payload['brokers'] = brokers
    client.upsert(collection_name='trusts', points=[{'id': tp.id, 'vector': embed_text(payload.get('text','')), 'payload': payload}])
    return SimpleResponse(ok=True, detail='broker added')

@router.get('/members')
def list_members():
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        return {'members': []}
    return {'members': tp.payload.get('members', [])}

@router.get('/brokers')
def list_brokers():
    client = get_client()
    tp = _get_trust_point(client)
    if not tp:
        return {'brokers': []}
    return {'brokers': tp.payload.get('brokers', [])}
