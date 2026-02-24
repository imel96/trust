from decimal import Decimal
from typing import List

def total_cash_flow_type(transactions: List, type: str) -> Decimal:
    totals = {'Investing': Decimal(0), 'Operating': Decimal(0), 'Financing': Decimal(0)}

    for tx in transactions:
        if tx['security'] is not None:
            if tx['units'] and tx['units'] != 0:
                totals['Investing'] += Decimal(tx['amount']) or Decimal("0")

    return totals[type]

def sum_account(transactions: List, account: str) -> Decimal:
    totals = dict()

    for tx in transactions:
        if tx['cash_account'] not in totals:
            totals[tx['cash_account']] = Decimal("0")
        if tx['offset_account'] not in totals:
            totals[tx['offset_account']] = Decimal("0")
        totals[tx['cash_account']] += Decimal(tx['amount']) or Decimal("0")
        if tx['offset_account'] is not None:
            totals[tx['offset_account']] -= Decimal(tx['amount']) or Decimal("0")
    if account == 'shares':
        print('accounting', transactions)
    if account in totals:
        return totals[account]
    else:
        return Decimal('0')

def security_balance(transactions: List, code: str) -> int:
    totals = dict()
    print(transactions)

    for tx in transactions:
        if tx['security'] not in totals:
            totals[tx['security']] = 0
        if tx['security'] is not None and tx['units'] is not None:
            totals[tx['security']] += tx['units']
    return totals[code]
