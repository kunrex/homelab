package internal

import (
	"encoding/json"
	"fmt"

	"mail-tui/llm"
	"mail-tui/models"
)

type mailSummary struct {
	Subject  string `json:"subject"`
	Sender   string `json:"sender"`
	Category string `json:"category"`
	Summary  string `json:"one_line"`
	Date     string `json:"date"`
}

const (
	SYSTEM_PROMPT = `You are a helpful assistant. Summarise the following emails into a clean Markdown digest.
					 Use ## headings per category, bullet points per mail.
					 For each mail include: sender, subject, and a one-line summary. Be concise.`
)

func GenerateDigest(client *llm.Client, mails []models.ClassifiedMail) (string, error) {
	summaries := make([]mailSummary, len(mails))
	for i, m := range mails {
		summaries[i] = mailSummary{
			Subject:  m.Mail.Subject,
			Sender:   m.Mail.Sender,
			Category: m.Category,
			Summary:  m.OneLine,
			Date:     m.Mail.Date,
		}
	}

	data, _ := json.Marshal(summaries)
	return client.Complete(SYSTEM_PROMPT, fmt.Sprintf("Emails:\n%s", string(data)))
}
