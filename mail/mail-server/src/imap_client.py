import email
import logging
import os
from datetime import datetime, timezone
from email.header import decode_header
from typing import List

import imapclient

from .datatypes import Config, MailRecord
from .utils import make_mail_id

logger = logging.getLogger(__name__)

def _decode_header_value(raw) -> str:
    parts = decode_header(raw or "")
    out = []
    for value, enc in parts:
        if isinstance(value, bytes):
            out.append(value.decode(enc or "utf-8", errors="replace"))
        else:
            out.append(str(value))
    return "".join(out)

def _decode_bytes(value) -> str:
    if isinstance(value, bytes):
        return value.decode("utf-8", errors="replace")
    return str(value) if value else ""

def _extract_body(msg: email.message.Message) -> tuple[str, str]:
    body_text, body_html = "", ""
    if msg.is_multipart():
        for part in msg.walk():
            ct = part.get_content_type()
            if ct == "text/plain" and not body_text:
                body_text = _decode_bytes(part.get_payload(decode=True))
            elif ct == "text/html" and not body_html:
                body_html = _decode_bytes(part.get_payload(decode=True))
    else:
        payload = _decode_bytes(msg.get_payload(decode=True))
        if msg.get_content_type() == "text/html":
            body_html = payload
        else:
            body_text = payload
    return body_text, body_html

class IMAPPoller:
    def __init__(self, cfg: Config):
        self.acct = cfg.account
        self.label = cfg.account.label
        self.password = os.environ.get("IMAP_PASSWORD", "")

    def _sender_allowed(self, raw_from: str) -> bool:
        if not self.acct.allowed_senders:
            return True
        addr = raw_from.lower()
        return any(allowed.lower() in addr for allowed in self.acct.allowed_senders)

    def poll(self) -> List[MailRecord]:
        mails: List[MailRecord] = []
        try:
            with imapclient.IMAPClient(
                self.acct.host, port=self.acct.port, ssl=True
            ) as client:
                client.login(self.acct.username, self.password)

                for folder in self.acct.folders:
                    client.select_folder(folder, readonly=True)
                    uids = client.search(["UNSEEN"])
                    if not uids:
                        continue

                    fetched = client.fetch(uids, ["RFC822"])
                    for uid, data in fetched.items():
                        try:
                            msg = email.message_from_bytes(data[b"RFC822"])
                            raw_from = msg.get("From", "")
                            if not self._sender_allowed(raw_from):
                                logger.debug("skipping uid %s from %s (not in allowed_senders)", uid, raw_from)
                                continue
                            body_text, body_html = _extract_body(msg)
                            mail_id = make_mail_id(self.label, str(uid))
                            mails.append(
                                MailRecord(
                                    id=mail_id,
                                    account=self.label,
                                    uid=str(uid),
                                    sender=raw_from,
                                    recipient=msg.get("To", ""),
                                    subject=_decode_header_value(msg.get("Subject", "")),
                                    date=msg.get("Date", ""),
                                    body_text=body_text,
                                    body_html=body_html,
                                    processed=False,
                                    ingested_at=datetime.now(timezone.utc).isoformat(),
                                )
                            )
                        except Exception as e:
                            logger.warning("skipping uid %s: %s", uid, e)
        except Exception as e:
            logger.error("IMAP poll failed: %s", e)
        return mails
