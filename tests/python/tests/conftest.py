from pytest_bdd import scenarios, given, when, then, parsers
import subprocess, time, os, requests, signal
import pytest
from utils import get_trust_payload

import sys
sys.path.append(os.path.dirname(os.path.dirname(__file__)))
from app.qdrant_client import get_client
from app.settings import settings

BASE = 'http://localhost:8000'

@pytest.fixture(scope='session', autouse=True)
def start_server():
    # start uvicorn server for the duration of the tests
    proc = subprocess.Popen(["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(2)
    yield
    proc.terminate()

@pytest.fixture(autouse=True)
def reset_qdrant():
    """
    Reset Qdrant collections before each test scenario.
    Ensures a clean state so tests don't affect each other.
    """
    client = get_client()

    # Collections you want to manage
    collections = ["trusts", "persons", "companies", "brokers"]

    for col in collections:
        try:
            client.delete_collection(col)
        except Exception:
            # Ignore if collection doesn't exist
            pass

        # Recreate collection fresh
        client.create_collection(
            collection_name=col,
            vectors_config={
                "size": settings.QDRANT_VECTOR_SIZE,
                "distance": "Cosine",
            },
        )
    yield

@given("there's a broker in the database")
def broker_exists():
    r = requests.post(f"{BASE}/entities/broker", json={"name": "Broker1"})
    assert r.status_code == 200

@given(parsers.parse('the trust has {n:d} shares in the trust portfolio'))
def trust_has_shares(n):
    payload = get_trust_payload()
    portfolio = payload.get('portfolio',{})
    portfolio['BHP'] = n
    payload['portfolio'] = portfolio
    r = requests.post(f"{BASE}/entities/trust", json=payload)
    assert r.status_code == 200

@given("there's a trust in the database")
def trust_exists():
    r = requests.post(f"{BASE}/entities/trust", json={"name": "Unit Trust"})
    assert r.status_code == 200

