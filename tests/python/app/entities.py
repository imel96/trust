from pydantic import BaseModel
from typing import Optional

class TrustIn(BaseModel):
    id: Optional[str] = None
    name: str
    portfolio: dict = {}
    cash: int = 0
    shares: int = 0

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
