# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/26 22:47
# -------------------------------------------
from fastapi import APIRouter

from .task import router as task_route
from .template import router as template_route
from .scheduled import router as scheduled_route
router = APIRouter()

router.include_router(task_route, prefix="/task")
router.include_router(template_route, prefix="/task/template")
router.include_router(scheduled_route, prefix="/task/scheduled")
