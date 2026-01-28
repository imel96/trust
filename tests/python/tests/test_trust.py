from pytest_bdd import scenarios, given, when, then
import requests
from utils import get_trust_payload

scenarios('../features/trust.feature')

BASE = 'http://localhost:8000'

@given("there's a person in the database")
def person_exists():
    r = requests.post(f"{BASE}/entities/person", json={"name": "Member1"})
    assert r.status_code == 200

@given("there's a broker in the database")
def broker_exists():
    r = requests.post(f"{BASE}/entities/broker", json={"name": "BrokerX"})
    assert r.status_code == 200

@when("the person is added to the trust")
def add_member():
    r = requests.post(f"{BASE}/trusts/add_member", params={'person_id': 'Member1'})
    assert r.status_code == 200

@then("the person can be seen in the trust member list")
def check_member():
    payload = get_trust_payload()
    assert "Member1" in payload.get("members", [])

@when("the broker is added to the trust")
def add_broker():
    r = requests.post(f"{BASE}/trusts/add_broker", params={'broker_id': 'BrokerX'})
    assert r.status_code == 200

@then("the broker can be seen in the trust broker list")
def check_broker():
    payload = get_trust_payload()
    assert "BrokerX" in payload.get("brokers", [])

