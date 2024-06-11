# -*- coding:utf-8 -*-　　
# @name: users
# @version:
from fastapi import APIRouter, Depends, HTTPException, Body
from pydantic import BaseModel
import jwt
from fastapi.security import OAuth2PasswordBearer
from datetime import datetime, timedelta
import hashlib
from core.config import SECRET_KEY
from core.db import get_mongo_db
router = APIRouter()

ALGORITHM = "HS256"

class LoginRequest(BaseModel):
    username: str
    password: str
class ChangePassword(BaseModel):
    newPassword: str
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


def create_access_token(data: dict, expires_delta: timedelta):
    to_encode = data.copy()
    expire = datetime.utcnow() + expires_delta
    to_encode.update({"exp": expire})
    encoded_jwt = jwt.encode(to_encode, SECRET_KEY, algorithm=ALGORITHM)
    return encoded_jwt


async def verify_token(token: str = Depends(oauth2_scheme)):
    credentials_exception = HTTPException(
        status_code=200,
        detail={"code": 401, "message": "Could not validate credentials"},
        headers={"WWW-Authenticate": "Bearer"},
    )
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        username: str = payload.get("sub")
        if username is None:
            raise credentials_exception
        return {"sub": username}
    except :
        raise credentials_exception


def hash_password(password: str) -> str:
    hashed_password = hashlib.sha256(password.encode()).hexdigest()
    return hashed_password


async def verify_user(username: str, password: str, db):
    user = await db.user.find_one({"username": username})
    if user and user["password"] == hash_password(password):
        return user
    return None


@router.post("/user/login")
async def login(login_request: LoginRequest = Body(...), db=Depends(get_mongo_db)):
    username = login_request.username
    password = login_request.password
    user = await verify_user(username, password, db)
    if user is None:
        return {
            'code': 401,
            'message' : 'Incorrect username or password'
          }

    token_data = {"sub": username}
    expires_delta = timedelta(days=30)  # Set the expiration time as needed
    token = create_access_token(token_data, expires_delta)

    return {
            'code': 200,
            'data': {
                'access_token': token
            }
          }


@router.post("/user/changePassword")
async def change_password(change_password: ChangePassword = Body(...), _: dict = Depends(verify_token),  db=Depends(get_mongo_db)):
    try:
        newPassword = hash_password(change_password.newPassword)
        await db.user.update_one({"username": 'ScopeSentry'},
                                  {"$set": {"password": newPassword}})
        return {
            'code': 200,
            'message': 'success change password'
        }
    except:
        return {
            'code': 500,
            'message': 'Password change failed'
        }