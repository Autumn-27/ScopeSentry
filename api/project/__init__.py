# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/29 20:58
# -------------------------------------------
from fastapi import APIRouter

from .project import router as project_route
router = APIRouter()

router.include_router(project_route, prefix="/project")
