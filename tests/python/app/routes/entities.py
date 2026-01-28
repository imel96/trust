from fastapi import APIRouter
from app.entities import TrustIn, PersonIn, BrokerIn, CompanyIn, SimpleResponse
from app.qdrant_client import get_client, ensure_collection
from app.rag_utils import embed_text
import uuid

router = APIRouter(prefix='/entities', tags=['entities'])

def _store(collection, payload):
    client = get_client()
    ensure_collection(collection)
    uid = payload.get('id') or str(uuid.uuid4())
    payload['id'] = uid
    vector = embed_text(payload.get('text',''))
    client.upsert(collection_name=collection, points=[{'id': uid, 'vector': vector, 'payload': payload}])
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    return uid

@router.post('/trust', response_model=SimpleResponse)
def create_trust(trust: TrustIn):
    payload = {'id': trust.id or None, 'name': trust.name, 'cash': trust.cash, 'income': 0.0, 'cost': 0.0, 'portfolio': trust.portfolio, 'members': [], 'brokers': [], 'units_outstanding':0, 'nav_per_unit':0.0, 'shares': trust.shares, 'text': f"trust:{trust.name}"}
    uid = _store('trusts', payload)
    return SimpleResponse(ok=True, detail=uid)

@router.post('/person', response_model=SimpleResponse)
def create_person(person: PersonIn):
    payload = {'id': person.id or None, 'name': person.name, 'text': f"person:{person.name}", 'units':0}
    uid = _store('persons', payload)
    return SimpleResponse(ok=True, detail=uid)

@router.post('/broker', response_model=SimpleResponse)
def create_broker(broker: BrokerIn):
    payload = {'id': broker.id or None, 'name': broker.name, 'text': f"broker:{broker.name}"}
    uid = _store('brokers', payload)
    return SimpleResponse(ok=True, detail=uid)

@router.post('/company', response_model=SimpleResponse)
def create_company(company: CompanyIn):
    payload = {'id': company.id or None, 'name': company.name, 'text': f"company:{company.name}"}
    uid = _store('companies', payload)
    return SimpleResponse(ok=True, detail=uid)
