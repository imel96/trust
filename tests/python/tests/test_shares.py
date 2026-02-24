from app.services.accounting import *
from decimal import Decimal
from pytest_bdd import scenarios, given, when, then, parsers
import requests
from utils import get_trust_payload

scenarios('../features/investing.feature')

BASE = 'http://localhost:8000'

@given(parsers.parse('the trust has {amount:d} dollar in "{account}" account'))
def set_account(amount, account):
    payload = get_trust_payload()
    payload[account] = float(amount)
    requests.post(f"{BASE}/entities/trust", json=payload)

@given(parsers.parse('the company "{code}" is on the market'))
def company_exists(code):
    r = requests.post(f"{BASE}/entities/company", json={"name": code})
    assert r.status_code == 200

@when(parsers.parse('the trust "{action}" {qty:d} shares of "{code}" through the broker'))
def trade(action, qty, code):
    if action == "buys":
        r = requests.post(f"{BASE}/shares/buy", params={"code": code, "qty": qty, "price": 10})
    else:
        r = requests.post(f"{BASE}/shares/sell", params={"code": code, "qty": qty})
    assert r.status_code == 200

@then(parsers.parse('"{code}" shares show up in the trust portfolio'))
def check_portfolio(code):
    payload = get_trust_payload()
    assert security_balance(payload['transactions'], code) != 0, f"{code} not in portfolio"

@then(parsers.parse('the trust portfolio shares of "{code}" becomes {qty:d}'))
def check_portfolio_qty(code, qty):
    payload = get_trust_payload()
    assert security_balance(payload['transactions'], code) == qty

@then(parsers.parse('the trust "{account}" account becomes {amount:d} dollar'))
def check_account(account, amount):
    payload = get_trust_payload()
    assert sum_account(payload['transactions'], account) == Decimal(amount)

@then(parsers.parse('the total "{cash_flow}" cash flow will be {amount:d} dollar'))
def check_account(cash_flow, amount):
    payload = get_trust_payload()
    assert total_cash_flow_type(payload['transactions'], cash_flow) == Decimal(amount)

