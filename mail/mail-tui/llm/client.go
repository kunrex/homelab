package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"mail-tui/cfg"
)

type Client struct {
	endpoint    string
	model       string
	apiKey      string
	temperature float64
	maxTokens   int
	topP        float64
	http        *http.Client
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
	TopP        *float64      `json:"top_p,omitempty"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

type modelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func NewClient(llmCfg cfg.LLMConfig) *Client {
	return &Client{
		endpoint:    llmCfg.Endpoint,
		model:       llmCfg.Model,
		apiKey:      llmCfg.APIKey,
		temperature: llmCfg.Temperature,
		maxTokens:   llmCfg.MaxTokens,
		topP:        llmCfg.TopP,
		http:        &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) Complete(system, user string) (string, error) {
	req := chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
	}
	if c.temperature > 0 {
		req.Temperature = &c.temperature
	}
	if c.maxTokens > 0 {
		req.MaxTokens = &c.maxTokens
	}
	if c.topP > 0 {
		req.TopP = &c.topP
	}

	body, _ := json.Marshal(req)
	httpClient := &http.Client{}
	httpReq, err := http.NewRequest("POST", c.endpoint+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("LLM request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("LLM returned %d", resp.StatusCode)
	}
	var cr chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		return "", fmt.Errorf("decode LLM response: %w", err)
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("empty choices from LLM")
	}
	return cr.Choices[0].Message.Content, nil
}

func (c *Client) Status() (up bool, latency time.Duration, err error) {
	start := time.Now()
	resp, err := c.http.Get(c.endpoint + "/v1/models")
	if err != nil {
		return false, 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, time.Since(start), nil
}

func (c *Client) Models() ([]string, error) {
	resp, err := c.http.Get(c.endpoint + "/v1/models")
	if err != nil {
		return nil, fmt.Errorf("GET /v1/models: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}
	var mr modelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&mr); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	ids := make([]string, len(mr.Data))
	for i, m := range mr.Data {
		ids[i] = m.ID
	}
	return ids, nil
}
