package Investigator

import "Investigator/helper/llm/vars"

type Model struct {
	APISecret    string `json:"APISecret,omitempty"`
	APIUrl       string `json:"APIUrl,omitempty"`
	APIModelName string `json:"APIModelName,omitempty"`
}

type ModelGraphState struct {
	Topic     string `json:"Topic,omitempty"`
	KimiModel *Model `json:"KimiModel,omitempty"`
	//DeepSeekModel *Model `json:"DeepSeekModel,omitempty"`
}

type ModelStateOutput struct {
	SEO      []string `json:"seo,omitempty"`
	Analysis string   `json:"analysis,omitempty"`
	Creation []struct {
		Inspiration string   `json:"inspiration,omitempty"`
		Content     string   `json:"content,omitempty"`
		Origin      []string `json:"origin,omitempty"`
	} `json:"creation,omitempty"`
	Response *vars.LLMResponse `json:"response,omitempty"`
}
