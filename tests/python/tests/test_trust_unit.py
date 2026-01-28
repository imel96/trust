from pytest_bdd import scenarios, given, when, then, parsers
import requests
from utils import get_trust_payload, get_person_payload

scenarios('../features/trust_unit.feature')

BASE = 'http://localhost:8000'

@given(parsers.parse('the trust has {amount:d} dollar in "{account}" account'))
def set_account(amount, account):
    payload = get_trust_payload()
    payload[account] = float(amount)
    requests.post(f"{BASE}/entities/trust", json=payload)

@given("there's a person in the trust")
def person_in_trust():
    requests.post(f"{BASE}/entities/person", json={"name": "InvestorUnit"})

@given(parsers.parse('the person owns {amount:d} unit in the trust'))
def person_units_amount(amount):
    payload = get_person_payload("InvestorUnit")
    payload['units'] = float(amount)
    requests.post(f"{BASE}/entities/person", json=payload)

@when(parsers.parse('the person "{action}" {amount:d} dollar'))
def person_action(action, amount):
    if action == "invest":
        r = requests.post(f"{BASE}/trust-unit/invest", params={'person_id': 'InvestorUnit', 'amount': amount})
    else:
        r = requests.post(f"{BASE}/trust-unit/redeem", params={'person_id': 'InvestorUnit', 'amount': amount})
    assert r.status_code == 200

@then(parsers.parse('the trust\'s "{account}" account become {amount:d} dollar'))
def check_account(account, amount):
    payload = get_trust_payload()
    assert round(payload[account], 2) == float(amount)

@then(parsers.parse('the person will have {amount:d} unit in the trust'))
def check_person_units(amount):
    person = get_person_payload("InvestorUnit")
    assert person["units"] == amount

@then(parsers.parse('the NAV per unit will be {amount:d} dollar'))
def check_nav(amount):
    payload = get_trust_payload()
    assert round(payload["nav_per_unit"], 2) == float(amount)

@then(parsers.parse('the trust will have {amount:d} unit outstanding'))
def check_units_outstanding(amount):
    payload = get_trust_payload()
    assert payload["units_outstanding"] == amount

