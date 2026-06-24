package models

type Mail struct {
	ID         string `json:"id"`
	Account    string `json:"account"`
	UID        string `json:"uid"`
	Sender     string `json:"sender"`
	Recipient  string `json:"recipient"`
	Subject    string `json:"subject"`
	Date       string `json:"date"`
	BodyText   string `json:"body_text"`
	BodyHTML   string `json:"body_html"`
	Processed  bool   `json:"processed"`
	IngestedAt string `json:"ingested_at"`
}

type ClassifiedMail struct {
	Mail     Mail
	Category string
	OneLine  string
}
