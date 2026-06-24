import json

from pathlib import Path
from typing import Optional
from datetime import datetime, timezone, timedelta

from .models import MailRecord

class MailStore:
    def __init__(self, data_dir: Path, ttl_days: int):
        self.data_dir = data_dir
        self.ttl_days = ttl_days
        self.data_dir.mkdir(parents=True, exist_ok=True)

    def _path(self, mail_id: str) -> Path:
        return self.data_dir / f"{mail_id}.json"

    def save(self, mail: MailRecord) -> bool:
        path = self._path(mail.id)
        if path.exists():
            return False
        path.write_text(json.dumps(mail.model_dump(), indent=2))
        return True

    def get(self, mail_id: str) -> Optional[MailRecord]:
        path = self._path(mail_id)
        if not path.exists():
            return None
        return MailRecord(**json.loads(path.read_text()))

    def list_unprocessed(self) -> list[MailRecord]:
        mails = []
        for f in self.data_dir.glob("*.json"):
            data = json.loads(f.read_text())
            if not data.get("processed", False):
                mails.append(MailRecord(**data))
        return sorted(mails, key=lambda m: m.ingested_at)

    def mark_processed(self, mail_id: str) -> bool:
        path = self._path(mail_id)
        if not path.exists():
            return False
        data = json.loads(path.read_text())
        data["processed"] = True
        path.write_text(json.dumps(data, indent=2))
        return True

    def mark_all_processed(self) -> int:
        count = 0
        for f in self.data_dir.glob("*.json"):
            data = json.loads(f.read_text())
            if not data.get("processed", False):
                data["processed"] = True
                f.write_text(json.dumps(data, indent=2))
                count += 1
        return count

    def delete(self, mail_id: str) -> bool:
        path = self._path(mail_id)
        if not path.exists():
            return False
        path.unlink()
        return True

    def delete_all(self) -> int:
        count = 0
        for f in self.data_dir.glob("*.json"):
            f.unlink()
            count += 1
        return count

    def cleanup_expired(self) -> int:
        cutoff = datetime.now(timezone.utc) - timedelta(days=self.ttl_days)
        count = 0
        for f in self.data_dir.glob("*.json"):
            data = json.loads(f.read_text())
            ingested = datetime.fromisoformat(data["ingested_at"])
            if ingested.tzinfo is None:
                ingested = ingested.replace(tzinfo=timezone.utc)
            if ingested < cutoff:
                f.unlink()
                count += 1
        return count
