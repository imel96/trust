from fastapi import FastAPI
from app.routes import entities, income, shares, trust, trust_unit, rag
from app.settings import settings

app = FastAPI(title='Trust App')

app.include_router(entities.router)
app.include_router(income.router)
app.include_router(shares.router)
app.include_router(trust.router)
app.include_router(trust_unit.router)
#app.include_router(rag.router)

@app.get('/')
def root():
    return {'message': 'Trust App API is running', 'ollama': settings.OLLAMA_URL}
