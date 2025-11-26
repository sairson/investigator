package Investigator

import (
	"Investigator/helper/llm"
	"Investigator/helper/llm/vars"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// MediaResearch 构建媒体调研的函数
func MediaResearch(ctx context.Context, state ModelGraphState) (vars.LLMResponse, error) {
	kimiClient := llm.NewKimiModel(ctx, state.KimiModel.APISecret, state.KimiModel.APIUrl, state.KimiModel.APIModelName)
	if kimiClient == nil {
		return vars.LLMResponse{}, fmt.Errorf("kimi chat model initialize nil")
	}
	content, err := kimiClient.Generate(
		[]map[string]interface{}{
			{"role": "system", "content": RoleSystem()},
			{"role": "user", "content": fmt.Sprintf("请你搜索全网查找有关话题,尽可能的理解相关话题的内容，总结并输出，我的话题或热点是: %s", state.Topic)},
		},
	)
	if err != nil {
		return vars.LLMResponse{}, fmt.Errorf("kimi chat model generate content failure %s", err)
	}
	return vars.LLMResponse{Content: content.Content, SearchContentTotalTokens: content.SearchContentTotalTokens, ChatPromptTokens: content.ChatPromptTokens, ChatCompletionTokens: content.ChatCompletionTokens, ChatTotalTokens: content.ChatTotalTokens}, nil
}

func InputFormatting(ctx context.Context, input vars.LLMResponse) (ModelStateOutput, error) {
	content := input.Content
	fmt.Println(content)
	content = strings.ReplaceAll(content, "```json", "")
	content = strings.ReplaceAll(content, "```", "")
	reComment := regexp.MustCompile(`(?m)^\s*//.*$`)
	content = reComment.ReplaceAllString(content, "")
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start == -1 || end == -1 || end <= start {
		return ModelStateOutput{}, fmt.Errorf("invalid content, json not found")
	}
	jsonStr := content[start : end+1]
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal([]byte(jsonStr), &envelope); err != nil {
		return ModelStateOutput{}, fmt.Errorf("json parse error: %v", err)
	}
	var analysis string
	if v, ok := envelope["Analysis"]; ok {
		var aStr string
		if err := json.Unmarshal(v, &aStr); err == nil {
			analysis = strings.TrimSpace(aStr)
		} else {
			analysis = strings.TrimSpace(string(v))
		}
	}
	var seo []string
	if v, ok := envelope["SEO"]; ok {
		var arr []string
		if err := json.Unmarshal(v, &arr); err == nil {
			for _, s := range arr {
				s = strings.TrimSpace(s)
				if s != "" {
					seo = append(seo, s)
				}
			}
		} else {
			// 兼容数组项为非字符串的情况
			var anyArr []interface{}
			if err2 := json.Unmarshal(v, &anyArr); err2 == nil {
				for _, it := range anyArr {
					s := strings.TrimSpace(fmt.Sprint(it))
					if s != "" {
						seo = append(seo, s)
					}
				}
			} else {
				var s string
				if err3 := json.Unmarshal(v, &s); err3 == nil {
					parts := strings.Split(s, ",")
					for _, p := range parts {
						ps := strings.TrimSpace(p)
						if ps != "" {
							seo = append(seo, ps)
						}
					}
				}
			}
		}
	}
	clean := func(s string) string {
		s = strings.TrimSpace(s)
		if len(s) >= 2 && s[0] == '`' && s[len(s)-1] == '`' {
			s = strings.TrimSpace(s[1 : len(s)-1])
		}
		return s
	}

	var creation []struct {
		Inspiration string   `json:"inspiration,omitempty"`
		Content     string   `json:"content,omitempty"`
		Origin      []string `json:"origin,omitempty"`
	}

	if v, ok := envelope["Creation"]; ok {
		var items []map[string]interface{}
		if err := json.Unmarshal(v, &items); err != nil {
			var single map[string]interface{}
			if err2 := json.Unmarshal(v, &single); err2 == nil {
				items = []map[string]interface{}{single}
			} else {
				return ModelStateOutput{}, fmt.Errorf("json parse error in Creation: %v", err)
			}
		}
		for _, item := range items {
			var inspect string
			for _, k := range []string{"创作灵感", "标题", "inspiration", "title", "Inspiration", "Title"} {
				if t, ok := item[k].(string); ok {
					inspect = strings.TrimSpace(t)
					if inspect != "" {
						break
					}
				}
			}
			var cont string
			for _, k := range []string{"内容", "content", "text", "Content", "Text", "description"} {
				if c, ok := item[k].(string); ok {
					cont = strings.TrimSpace(c)
					if cont != "" {
						break
					}
				}
			}
			var origin []string
			found := false
			for _, k := range []string{"来源", "sources", "source", "origin", "url", "Origin", "Sources", "ORIGIN"} {
				if arr, ok := item[k].([]interface{}); ok {
					for _, v2 := range arr {
						s := strings.TrimSpace(fmt.Sprint(v2))
						s = clean(s)
						if s != "" {
							origin = append(origin, s)
						}
					}
					if len(origin) > 0 {
						found = true
						break
					}
				} else if s1, ok := item[k].(string); ok {
					parts := strings.Split(s1, ",")
					for _, p := range parts {
						ps := clean(p)
						if ps != "" {
							origin = append(origin, ps)
						}
					}
					if len(origin) > 0 {
						found = true
						break
					}
				}
			}
			_ = found
			if inspect == "" && cont == "" {
				b, _ := json.Marshal(item)
				cont = string(b)
			}
			creation = append(creation, struct {
				Inspiration string   `json:"inspiration,omitempty"`
				Content     string   `json:"content,omitempty"`
				Origin      []string `json:"origin,omitempty"`
			}{Inspiration: inspect, Content: cont, Origin: origin})
		}
	}
	return ModelStateOutput{SEO: seo, Analysis: analysis, Creation: creation, Response: &input}, nil
}
