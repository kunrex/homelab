from datetime import datetime, timezone

from apscheduler.schedulers.asyncio import AsyncIOScheduler

def init_scheduler(scheduler: AsyncIOScheduler, poll_fn, cleanup_fn, poll_interval_seconds: int):
    scheduler.add_job(
        poll_fn,
        "interval",
        seconds=poll_interval_seconds,
        id="imap_poll",
        max_instances=1,
        coalesce=True,
        next_run_time=datetime.now(timezone.utc),
    )
    scheduler.add_job(
        cleanup_fn,
        "interval",
        hours=1,
        id="ttl_cleanup",
        max_instances=1,
        coalesce=True,
    )
