import os
import asyncio
import logging

from pathlib import Path
from contextlib import asynccontextmanager

from fastapi import FastAPI, APIRouter
from apscheduler.schedulers.asyncio import AsyncIOScheduler

from src import IMAPPoller, MailStore, load_config, init_routes, init_scheduler

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
logger = logging.getLogger(__name__)

cfg = load_config()
if cfg is None:
    print("failed to validate config")
    exit(1)

router = APIRouter()

store = MailStore(Path(os.getenv("DATA_DIR","")), cfg.mail_ttl_days)
init_routes(router, store)

poller = IMAPPoller(cfg)

def poll_and_store():
    mails = poller.poll()
    new = sum(1 for m in mails if store.save(m))
    if new:
        logger.info("Ingested %d new mail(s)", new)

def cleanup():
    removed = store.cleanup_expired()
    if removed:
        logger.info("Expired %d mail(s)", removed)

@asynccontextmanager
async def lifespan(app: FastAPI):
    if cfg is not None:
        scheduler = AsyncIOScheduler()
        init_scheduler(scheduler, poll_and_store, cleanup, cfg.poll_interval_seconds)
        scheduler.start()

        asyncio.get_event_loop().run_in_executor(None, poll_and_store)

        yield

        scheduler.shutdown(wait=False)

app = FastAPI(title="mail-server", lifespan=lifespan)
app.include_router(router)

@app.get("/health")
def health():
    return {"ok": True}
