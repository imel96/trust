import requests
import os

QDRANT_URL = os.getenv("QDRANT_URL", "http://qdrant:6333")

def get_trust_payload():
    """Get the first (and only) trust payload from Qdrant."""
    res = requests.post(
        f"{QDRANT_URL}/collections/trusts/points/scroll",
        json={"limit": 1}
    ).json()
    points = res.get("result", {}).get("points", [])
    if not points:
        raise AssertionError("No trust found in Qdrant")
    return points[0]["payload"]

def get_person_payload(name: str):
    """Find a person by name."""
    res = requests.post(
        f"{QDRANT_URL}/collections/persons/points/scroll",
        json={"limit": 100}
    ).json()
    for p in res.get("result", {}).get("points", []):
        if p["payload"].get("name") == name:
            return p["payload"]
    raise AssertionError(f"Person '{name}' not found in Qdrant")

