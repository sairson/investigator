package vars

type ToolDefine struct {
	Type     string                 `json:"type"`
	Function map[string]interface{} `json:"function"`
}

type ChatRequest struct {
	Tools          []ToolDefine             `json:"tools"`
	Model          string                   `json:"model"`
	MaxTokens      int                      `json:"max_tokens"`
	Messages       []map[string]interface{} `json:"messages"`
	Temperature    float64                  `json:"temperature"`
	ResponseFormat map[string]string        `json:"response_format"`
}

type ToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls"`
}

type Choice struct {
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type ChatResponse struct {
	Choices []Choice   `json:"choices"`
	Usage   TokenUsage `json:"usage"`
}

type LLMResponse struct {
	Content                  string `json:"content"`
	SearchContentTotalTokens int    `json:"search_content_total_tokens"`
	ChatPromptTokens         int    `json:"chat_prompt_tokens"`
	ChatCompletionTokens     int    `json:"chat_completion_tokens"`
	ChatTotalTokens          int    `json:"chat_total_tokens"`
}
