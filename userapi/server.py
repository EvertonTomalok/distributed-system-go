from typing import Union

from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


class UserResponse(BaseModel):
    user_id: str
    is_valid: bool


@app.get("/")
def read_root():
    return {"status": "ok"}


@app.get("/api/user/{user_id}", response_model=UserResponse)
def read_item(user_id: str):
    is_valid = "invalid" not in user_id
    return {"user_id": user_id, "is_valid": is_valid}
