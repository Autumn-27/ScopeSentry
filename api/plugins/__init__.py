# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/24 19:54
# -------------------------------------------
from fastapi import APIRouter

from .plugin import router as plugin_routeer
router = APIRouter()

router.include_router(plugin_routeer, prefix="/plugin")
