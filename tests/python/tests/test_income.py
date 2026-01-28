from pytest_bdd import scenarios, given, when, then, parsers
import requests
from utils import get_trust_payload

scenarios('../features/income.feature')

BASE = 'http://localhost:8000'

@given("there's a trust in the database")
def trust_exists():
    r = requests.post(f"{BASE}/entities/trust", json={"name": "Test Trust"})
    assert r.status_code == 200

@given(parsers.parse('the trust has {amount:d} dollar in "{account}" account'))
def set_account(amount, account):
    payload = get_trust_payload()
    payload[account] = float(amount)
    requests.post(f"{BASE}/entities/trust", json=payload)

@given("there's a company in the database")
def company_exists():
    r = requests.post(f"{BASE}/entities/company", json={"name": "BHP"})
    assert r.status_code == 200

@when(parsers.parse('the company pays dividend of {dollar:g} dollar per share'))
def pay_dividend(dollar):
    r = requests.post(f"{BASE}/income/dividend", params={'dollar_per_share': dollar})
    assert r.status_code == 200

@then(parsers.parse('the trust\'s "{account}" account become {amount:d} dollar'))
def check_account(account, amount):
    payload = get_trust_payload()
    assert round(payload[account], 2) == float(amount), \
        f"Expected {amount}, got {payload[account]}"

