package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"mail-tui/cfg"
	"mail-tui/llm"
	"mail-tui/models"
)

type classifyResult struct {
	Relevant bool   `json:"relevant"`
	Category string `json:"category"`
	OneLine  string `json:"one_line"`
}

const (
	SYSTEM_PROMT = `You are a filter for a student's inbox. Given an email, respond ONLY with a SINGLE valid JSON OBJECT and DO NOT OVERGENERATE:
					{"relevant": true/false, "category": "<category or null>", "one_line": "<short reason>"}
					Relevant means it relates to: %s. DONT include any other form of reasoning, ONLY ONE VALID JSON OBJECT`
)

func ClassifyMails(client *llm.Client, mails []models.Mail) ([]models.ClassifiedMail, error) {
	system := fmt.Sprintf(SYSTEM_PROMT, cfg.Cfg.CategoriesStr)

	results := make([]models.ClassifiedMail, len(mails))
	errs := make([]error, len(mails))
	var wg sync.WaitGroup
	var done atomic.Int64

	for i, mail := range mails {
		wg.Add(1)
		go func(idx int, m models.Mail) {
			defer wg.Done()
			defer func() {
				n := done.Add(1)
				fmt.Printf("\r\033[K  %s  %s / %s",
					stepStyle.Render("Classifying"),
					countStyle.Render(fmt.Sprintf("%d", n)),
					dimStyle.Render(fmt.Sprintf("%d", len(mails))),
				)
			}()

			body := m.BodyText
			if len(body) > 1000 {
				body = body[:1000]
			}
			prompt := fmt.Sprintf("Subject: %s\nBody:\n%s", m.Subject, body)

			raw, err := client.Complete(system, prompt)
			if err != nil {
				errs[idx] = err
				return
			}
			raw = strings.TrimSpace(raw)
			raw = strings.TrimPrefix(raw, "```json")
			raw = strings.TrimPrefix(raw, "```")
			raw = strings.TrimSuffix(raw, "```")
			raw = strings.TrimSpace(raw)

			var r classifyResult
			if err := json.Unmarshal([]byte(raw), &r); err != nil {
				errs[idx] = fmt.Errorf("parse classification for %s: %w", m.ID, err)
				return
			}
			if r.Relevant {
				results[idx] = models.ClassifiedMail{Mail: m, Category: r.Category, OneLine: r.OneLine}
			}
		}(i, mail)
	}
	wg.Wait()
	fmt.Print("\r\033[K") // clear progress line; caller prints the done step

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	var relevant []models.ClassifiedMail
	for _, r := range results {
		if r.Mail.ID != "" {
			relevant = append(relevant, r)
		}
	}
	return relevant, nil
}
