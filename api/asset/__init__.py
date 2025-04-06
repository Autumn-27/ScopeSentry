# -------------------------------------
# @file      : __init__.py
# @author    : Autumn
# @contact   : rainy-autumn@outlook.com
# @time      : 2024/10/20 19:12
# -------------------------------------------
from fastapi import APIRouter

from .asset import router as asset_routeer
from .subdomain import router as subdomain_route
from .url import router as url_route
from .crawler import router as crawler_route
from .common import router as common_route
from .sensitive import router as sens_route
from .page_monitoring import router as page_monitoring_route
from .dirscan import router as dirscan_route
from .SubdoaminTaker import router as subdoamintaker_route
from .vulnerability import router as vulnerability_route
from .statistics import router as statistics_route
from .root_domain import router as root_domain_route
from .app import router as app_route
from .mp import router as mp_route
router = APIRouter()

router.include_router(asset_routeer, prefix="/asset")
router.include_router(statistics_route, prefix="/asset/statistics")
router.include_router(subdomain_route, prefix="/subdomain")
router.include_router(url_route, prefix="/url")
router.include_router(crawler_route, prefix="/crawler")
router.include_router(common_route, prefix="/data")
router.include_router(sens_route, prefix="/sensitive")
router.include_router(page_monitoring_route, prefix="/page/monitoring")
router.include_router(dirscan_route, prefix="/dirscan/result")
router.include_router(subdoamintaker_route, prefix="/subdomaintaker")
router.include_router(vulnerability_route, prefix="/vul")
router.include_router(root_domain_route, prefix="/root/domain")
router.include_router(app_route, prefix="/app")
router.include_router(mp_route, prefix="/mp")

