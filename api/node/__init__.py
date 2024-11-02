# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/11/2 14:51
# -------------------------------------------
from fastapi import APIRouter

from .node import router as node_routeer
router = APIRouter()

router.include_router(node_routeer, prefix="/node")