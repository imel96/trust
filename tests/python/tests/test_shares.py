from pytest_bdd import scenarios, given, when, then, parsers
import requests
from utils import get_trust_payload

scenarios('../features/shares.feature')

BASE = 'http://localhost:8000'

@given(parsers.parse('the trust has {amount:d} dollar in "{account}" account'))
def set_account(amount, account):
    payload = get_trust_payload()
    payload[account] = float(amount)
    print(account, amount, payload)
    requests.post(f"{BASE}/entities/trust", json=payload)

@given("there's a company in the database")
def company_exists():
    r = requests.post(f"{BASE}/entities/company", json={"name": "BHP"})
    assert r.status_code == 200

@when(parsers.parse('the trust "{action}" {qty:d} shares of "{symbol}" through the broker'))
def trade(action, qty, symbol):
    if action == "buys":
        r = requests.post(f"{BASE}/shares/buy", params={"symbol": symbol, "qty": qty, "price": 10})
        #payload = get_trust_payload()
        #payload["cash"] -= 100 # credit
        #payload["cost"] += 1
        #requests.post(f"{BASE}/entities/trust", json=payload)
    else:
        r = requests.post(f"{BASE}/shares/sell", params={"symbol": symbol, "qty": qty})
    assert r.status_code == 200

@then(parsers.parse('"{symbol}" shares show up in the trust portfolio'))
def check_portfolio(symbol):
    payload = get_trust_payload()
    print(payload)
    assert symbol in payload["portfolio"], f"{symbol} not in portfolio"

@then(parsers.parse('the trust portfolio shares of "{symbol}" becomes {qty:d}'))
def check_portfolio_qty(symbol, qty):
    payload = get_trust_payload()
    assert payload["portfolio"].get(symbol, 0) == qty

@then(parsers.parse('the trust "{account}" account becomes {amount:d} dollar'))
def check_account(account, amount):
    payload = get_trust_payload()
    assert round(payload[account], 2) == float(amount)

