package chatgpt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	client          = http.Client{}
	URL             = "https://api.openai.com/v1/chat/completions"
	ChatGpt3_5Turbo = "gpt-3.5-turbo"
)

type ChatGptPrompt struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChapGptRequest struct {
	Model    string          `json:"model"`  //model name
	Stream   bool            `json:"stream"` // Response streaming
	Messages []ChatGptPrompt `json:"messages"`
}

type ChatResponse struct {
	ID        string               `json:"id"`
	Object    string               `json:"object"`
	CreatedAt int64                `json:"created_at"`
	Choices   []ChatResponseChoice `json:"choices"`
}

type ChatResponseChoice struct {
	Index        int           `json:"index"`
	Message      ChatGptPrompt `json:"delta"`
	FinishReason string        `json:"finish_reason"`
}

type ChatGptClient struct {
	Token    string
	Model    string
	Messages []ChatGptPrompt
}

type SSEEvent struct {
	Type string
	Data string
}

func NewChatGptClient(token string, options ...Option) *ChatGptClient {
	gpt := &ChatGptClient{Token: token, Model: ChatGpt3_5Turbo, Messages: make([]ChatGptPrompt, 0)}
	for _, option := range options {
		option(gpt)
	}
	return gpt
}
func (c *ChatGptClient) AddPrompt(role string, content string) {
	c.Messages = append(c.Messages, ChatGptPrompt{Role: role, Content: content})
}

func (c *ChatGptClient) SetPrompts(prompts []ChatGptPrompt) {
	c.Messages = prompts
}

func (c *ChatGptClient) SendRequestWithStream(role string, content string, closing <-chan time.Time) (chan ChatResponse, error) {
	c.Messages = append(c.Messages, ChatGptPrompt{Role: role, Content: content})

	body := ChapGptRequest{Model: c.Model, Stream: true, Messages: c.Messages}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", URL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("Authorization", "Bearer "+c.Token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Request failed")
	}
	out := make(chan ChatResponse)
	fmt.Println("hi")
	go func() {
		defer resp.Body.Close()
		defer close(out)
		fmt.Println("in routine")
		reader := bufio.NewScanner(resp.Body)
		for reader.Scan() {
			select {
			case <-closing:
				break
			default:
			}
			event, err := parseSSEEvent(reader.Text())
			if err != nil {
				fmt.Println(err)
				continue
			}
			if event.Data == "" {
				continue
			}
			var chat ChatResponse
			if err := json.Unmarshal([]byte(event.Data), &chat); err != nil {
				fmt.Println(err)
				break
			}
			out <- chat
		}
	}()
	fmt.Println("fire")
	return out, nil
}

func parseSSEEvent(line string) (SSEEvent, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return SSEEvent{}, fmt.Errorf("Invalid SSE event format: %s", line)
	}

	if strings.TrimSpace(parts[1]) == "[DONE]" {
		return SSEEvent{}, nil
	}
	event := SSEEvent{
		Type: strings.TrimSpace(parts[0]),
		Data: strings.TrimSpace(parts[1]),
	}

	return event, nil
}

func processSSEEvent(event SSEEvent) {
	fmt.Println("Event:", event)
}
