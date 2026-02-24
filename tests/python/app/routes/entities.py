from app.services.accounting import *
from datetime import datetime
from decimal import Decimal
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
    return uid

def _get_trust_point(client):
    ensure_collection('trusts')
    pts, _ = client.scroll(collection_name='trusts', limit=1)
    if pts:
        return pts[0]
    return None

@router.post('/trust', response_model=SimpleResponse)
def create_trust(trust: TrustIn):
    tp = _get_trust_point(get_client())
    transactions = []

    if tp:
        transactions = tp.payload['transactions']
    if trust.cash != 0:
        transactions = create_transactions_payload(transactions, trust.cash, 'cash', None, 0)
    if trust.shares != 0:
        transactions = create_transactions_payload(transactions, trust.shares, 'shares', None, 0)
    if trust.portfolio != {}:
        for security, asset in trust.portfolio.items():
            transactions = create_transactions_payload(transactions, 0, None, security, asset['units'])
    payload = {'id': trust.id or None, 'name': trust.name, 'income': 0.0, 'cost': 0.0, 'portfolio': trust.portfolio, 'members': [], 'brokers': [], 'text': f"trust:{trust.name}", 'transactions': transactions}
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

def create_transactions_payload(transactions: list, amount: float, cash_account: str, security: str, units: int) -> list:
    new_amount = Decimal(amount)
    if new_amount != Decimal(0):
        new_amount = new_amount + sum_account(transactions, cash_account)
    transactions.append({
        'transaction_date': datetime.now(),
        'security': security,
        'units': units,
        'price_per_unit': None,
        'amount': new_amount,
        'cash_account': cash_account,
        'offset_account': None,
    })
    return transactions
