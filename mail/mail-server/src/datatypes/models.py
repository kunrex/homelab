from pydantic import BaseModel, Field

class Account(BaseModel):
    label: str
    host: str
    port: int = 993
    username: str
    folders: list[str] = Field(default_factory=lambda: ["INBOX"])
    allowed_senders: list[str] = Field(default_factory=list)

class Config(BaseModel):
    account: Account
    mail_ttl_days: int
    poll_interval_seconds: int

class MailRecord(BaseModel):
    id: str
    account: str
    uid: str
    sender: str
    recipient: str
    subject: str
    date: str
    body_text: str
    body_html: str = ""
    processed: bool = False
    ingested_at: str

