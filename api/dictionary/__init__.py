# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/9/8 15:10
# -------------------------------------------
from fastapi import APIRouter

from .dir import router as dir_router
from .subdomain import router as subdoamin_router
router = APIRouter()

router.include_router(dir_router, prefix="/dir")
router.include_router(subdoamin_router, prefix="/subdomain")