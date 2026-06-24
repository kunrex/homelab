import hashlib


def make_mail_id(account: str, uid: str) -> str:
    return hashlib.sha256(f"{account}:{uid}".encode()).hexdigest()[:16]
