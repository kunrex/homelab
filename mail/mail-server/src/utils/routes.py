from typing import List

from fastapi import APIRouter, HTTPException

from ..datatypes import MailRecord, MailStore

def init_routes(router: APIRouter, store: MailStore):
    @router.get("/mails", response_model=List[MailRecord])
    def list_mails():
        if store is not None:
            return store.list_unprocessed()

    @router.get("/mails/{mail_id}", response_model=MailRecord)
    def get_mail(mail_id: str):
        if store is not None:
            mail = store.get(mail_id)
            if not mail:
                raise HTTPException(status_code=404, detail="Not found")

            return mail

    @router.post("/mails/process-all")
    def process_all():
        if store is not None:
            count = store.mark_all_processed()
            return {"processed": count}

    @router.post("/mails/{mail_id}/process")
    def process_mail(mail_id: str):
        if store is not None:
            if not store.mark_processed(mail_id):
                raise HTTPException(status_code=404, detail="Not found")

            return {"ok": True}

    @router.delete("/mails/{mail_id}")
    def delete_mail(mail_id: str):
        if store is not None:
            if not store.delete(mail_id):
                raise HTTPException(status_code=404, detail="Not found")

            return {"ok": True}

    @router.delete("/mails")
    def delete_all():
        if store is not None:
            count = store.delete_all()
            return {"deleted": count}
