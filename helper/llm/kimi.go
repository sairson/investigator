package llm

import (
	"Investigator/helper/llm/vars"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
)

type KimiModel struct {
	baseUrl    string
	model      string
	ChatModel  *openai.ChatModel
	apiKey     string
	httpClient *http.Client
	context    context.Context
}

func NewKimiModel(ctx context.Context, apiKey string, apiUrl string, modelName string) *KimiModel {
	return &KimiModel{
		baseUrl:    apiUrl,
		model:      modelName,
		apiKey:     apiKey,
		context:    ctx,
		httpClient: &http.Client{Timeout: 6000 * time.Second},
	}
}

func (k *KimiModel) Generate(messages []map[string]interface{}) (*vars.LLMResponse, error) {
	WebSearchTools := []vars.ToolDefine{
		{Type: "builtin_function", Function: map[string]interface{}{"name": "$web_search"}},
	}
	llmResponse := &vars.LLMResponse{}
	var finalContent string
	for {
		request := vars.ChatRequest{
			Model:          k.model,
			Messages:       messages,
			Temperature:    0.6,
			MaxTokens:      32768,
			Tools:          WebSearchTools,
			ResponseFormat: map[string]string{"type": "json_object"},
		}
		data, _ := json.Marshal(request)
		httpRequest, _ := http.NewRequestWithContext(k.context, http.MethodPost, k.baseUrl+"/chat/completions", bytes.NewReader(data))
		httpRequest.Header.Set("Content-Type", "application/json")
		httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", k.apiKey))
		httpResponse, err := k.httpClient.Do(httpRequest)
		if err != nil {
			return llmResponse, fmt.Errorf("llm kimi api invoke failed: %s", err.Error())
		}
		defer func() {
			_ = httpResponse.Body.Close()
		}()
		var chatResponse vars.ChatResponse
		if err := json.NewDecoder(httpResponse.Body).Decode(&chatResponse); err != nil {
			return llmResponse, fmt.Errorf("llm kimi api invoke failed: %s", err.Error())
		}
		if len(chatResponse.Choices) == 0 {
			return llmResponse, fmt.Errorf("llm kimi api response return empty choices")
		}
		choice := chatResponse.Choices[0]
		finish := choice.FinishReason
		messages = append(messages, map[string]interface{}{
			"role":       "assistant",
			"content":    choice.Message.Content,
			"tool_calls": choice.Message.ToolCalls,
		})
		if finish != "tool_calls" || len(choice.Message.ToolCalls) == 0 {
			//fmt.Printf("chat_prompt_tokens:          %d\n", chatResponse.Usage.PromptTokens)
			//fmt.Printf("chat_completion_tokens:      %d\n", chatResponse.Usage.CompletionTokens)
			//fmt.Printf("chat_total_tokens:           %d\n", chatResponse.Usage.TotalTokens)
			llmResponse.ChatPromptTokens = chatResponse.Usage.PromptTokens
			llmResponse.ChatCompletionTokens = chatResponse.Usage.CompletionTokens
			llmResponse.ChatTotalTokens = chatResponse.Usage.TotalTokens
			finalContent = choice.Message.Content
			break
		}
		for _, toolCall := range choice.Message.ToolCalls {
			if toolCall.Function.Name == "$web_search" {
				var args map[string]any
				_ = json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
				if usage, ok := args["usage"].(map[string]any); ok {
					if tt, ok := usage["total_tokens"]; ok {
						//fmt.Printf("search_content_total_tokens: %v\n", tt)
						llmResponse.SearchContentTotalTokens = int(tt.(float64))
					}
				}
				toolResult, _ := json.Marshal(args)
				messages = append(messages, map[string]interface{}{
					"role":         "tool",
					"tool_call_id": toolCall.ID,
					"name":         toolCall.Function.Name,
					"content":      string(toolResult),
				})
			} else {
				messages = append(messages, map[string]interface{}{
					"role":         "tool",
					"tool_call_id": toolCall.ID,
					"name":         toolCall.Function.Name,
					"content":      fmt.Sprintf("Error: unable to find tool by name '%s'", toolCall.Function.Name),
				})
			}
		}
	}
	llmResponse.Content = finalContent
	return llmResponse, nil
}
