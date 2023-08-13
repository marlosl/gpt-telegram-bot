package chatgpt

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model            string        `json:"model"`
	Messages         []ChatMessage `json:"messages"`
	Temperature      *int          `json:"temperature,omitempty"`
	TopP             *int          `json:"top_p,omitempty"`
	N                *int          `json:"n,omitempty"`
	Stream           *bool         `json:"stream,omitempty"`
	Stop             *string     `json:"stop,omitempty"`
	MaxTokens        *int          `json:"max_tokens,omitempty"`
	PresencePenalty  *float32      `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float32      `json:"frequency_penalty,omitempty"`
	User             *string       `json:"user,omitempty"`
}

type EditRequest struct {
	Model       string `json:"model"`
	Input       string `json:"input"`
	Instruction string `json:"instruction"`
	N           *int   `json:"n,omitempty"`
	Temperature *int   `json:"temperature,omitempty"`
	TopP        *int   `json:"top_p,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason *string     `json:"finish_reason,omitempty"`
}

type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type CreateImageRequest struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type UrlReponse struct {
	Url string `json:"url"`
}

type CreateImageResponse struct {
	Created int64        `json:"created"`
	Data    []UrlReponse `json:"data"`
}
