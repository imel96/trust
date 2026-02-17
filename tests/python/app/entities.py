from datetime import datetime
from decimal import Decimal
from pydantic import BaseModel, Field
from typing import List, Optional

class Transaction(BaseModel):
    id: Optional[str] = None
    transaction_date: datetime
    description: Optional[str] = None
    security: Optional[str] = None
    units: Optional[Decimal] = Decimal("0")
    price_per_unit: Optional[Decimal] = Decimal("0")
    amount: Optional[Decimal] = Decimal("0")
    franking_credit: Optional[Decimal] = Decimal("0")
    cash_flow: Optional[str] = None

class TrustIn(BaseModel):
    id: Optional[str] = None
    name: str
    portfolio: dict = {}
    cash: int = 0
    shares: int = 0
    transactions: list[Transaction] = Field(default_factory=list)

class PersonIn(BaseModel):
    id: Optional[str] = None
    name: str

class BrokerIn(BaseModel):
    id: Optional[str] = None
    name: str

class CompanyIn(BaseModel):
    id: Optional[str] = None
    name: str

class SimpleResponse(BaseModel):
    ok: bool
    detail: Optional[str] = None
