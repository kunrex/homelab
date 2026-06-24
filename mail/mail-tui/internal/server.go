package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"mail-tui/cfg"
	"mail-tui/models"
)

type ServerClient struct {
	base string
	http *http.Client
}

func NewServerClient() *ServerClient {
	return &ServerClient{
		base: fmt.Sprintf("http://%s:%d", cfg.Cfg.Server.Host, cfg.Cfg.Server.Port),
		http: &http.Client{},
	}
}

func (c *ServerClient) do(method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.base+path, nil)
	if err != nil {
		return nil, err
	}
	return c.http.Do(req)
}

func (c *ServerClient) FetchMails() ([]models.Mail, error) {
	resp, err := c.do("GET", "/mails")
	if err != nil {
		return nil, fmt.Errorf("GET /mails: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}
	var mails []models.Mail
	if err := json.NewDecoder(resp.Body).Decode(&mails); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return mails, nil
}

func (c *ServerClient) MarkAllProcessed() error {
	resp, err := c.do("POST", "/mails/process-all")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}
	return nil
}

func PromptMarkProcessed(server *ServerClient) {
	if !cfg.Cfg.Server.MarkProcessedAfterRun {
		return
	}
	if Confirm("Mark all mails as processed?") {
		if err := server.MarkAllProcessed(); err != nil {
			fmt.Fprintf(os.Stderr, "  warning: %v\n", err)
		} else {
			PrintStep("Done", "")
		}
	}
}
