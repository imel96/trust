from decimal import Decimal
from typing import List

def total_cash_flow_type(transactions: List, type: str) -> Decimal:
    totals = {'Investing': Decimal(0), 'Operating': Decimal(0), 'Financing': Decimal(0)}

    for tx in transactions:
        if  tx['security'] is not None:
            if tx['units'] and tx['units'] < 0:
                totals['Investing'] += Decimal(tx['amount']) or Decimal("0")
            if tx['units'] and tx['units'] > 0:
                totals['Investing'] -= Decimal(tx['amount']) or Decimal("0")

    return totals[type]
