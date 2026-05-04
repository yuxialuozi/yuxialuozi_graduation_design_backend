package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type AIService struct {
	apiKey string
}

func NewAIService() *AIService {
	return &AIService{
		apiKey: os.Getenv("BIGMODEL_API_KEY"),
	}
}

// AIMessage represents a chat message
type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AIChatRequest represents the chat request from frontend
type AIChatRequest struct {
	Messages []AIMessage `json:"messages"`
}

// AIStreamChunk represents a streaming response chunk
type AIStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

const (
	BigModelBaseURL = "https://open.bigmodel.cn/api/paas/v4"
	ModelName       = "glm-4-flash"
)

// System prompt for admin users
func getAdminSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的 AI 智能助手，专门帮助管理员进行数据分析和业务决策。

你的职责包括：
1. 数据分析：分析租户、合同、房间、费用、维修等数据，提供洞察和建议
2. 业务咨询：回答关于租金定价、合同条款、费用收取等业务问题
3. 趋势预测：基于历史数据，分析收入趋势、租户流失风险等
4. 维修建议：根据维修工单数据，提供维护优先级建议
5. 报表解读：帮助理解各类统计报表中的数据含义

请用专业、简洁的语言回答，如果有具体数据可以给出量化分析。
当前系统时间：%s`, time.Now().Format("2006-01-02 15:04:05"))
}

// System prompt for tenant (user) users
func getTenantSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的 AI 智能助手，专门帮助租户（住户）解决日常问题。

你的职责包括：
1. 费用咨询：解答关于租金、水电费、物业费等费用问题
2. 合同解读：帮助理解租房合同的条款和权益
3. 维修指引：指导如何提交维修申请，跟踪维修进度
4. 政策说明：解释租金调整、押金退还、续租等政策
5. 常见问题：回答租户常见的各类问题

请用友好、耐心的语言回答，使用通俗易懂的语言解释专业术语。
当前系统时间：%s`, time.Now().Format("2006-01-02 15:04:05"))
}

// ChatResult represents the result of a chat completion
type ChatResult struct {
	Content string
	Error   string
}

// Chat performs a non-streaming chat completion
func (s *AIService) Chat(role string, messages []AIMessage) (*ChatResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("AI API key not configured, please set BIGMODEL_API_KEY environment variable")
	}

	// Add system prompt based on role
	systemPrompt := getTenantSystemPrompt()
	if role == "admin" {
		systemPrompt = getAdminSystemPrompt()
	}

	// Prepare messages with system prompt
	allMessages := make([]AIMessage, 0, len(messages)+1)
	allMessages = append(allMessages, AIMessage{
		Role:    "system",
		Content: systemPrompt,
	})
	allMessages = append(allMessages, messages...)

	requestBody := map[string]interface{}{
		"model":    ModelName,
		"messages": allMessages,
		"stream":   false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", BigModelBaseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &ChatResult{Error: fmt.Sprintf("API error: %s", string(body))}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return nil, fmt.Errorf("invalid response format")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid choice format")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid message format")
	}

	content, ok := message["content"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid content format")
	}

	return &ChatResult{Content: content}, nil
}

// ChatStream performs a streaming chat completion
func (s *AIService) ChatStream(role string, messages []AIMessage) (<-chan string, <-chan error) {
	resultChan := make(chan string)
	errorChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		defer close(errorChan)

		if s.apiKey == "" {
			errorChan <- fmt.Errorf("AI API key not configured, please set BIGMODEL_API_KEY environment variable")
			return
		}

		// Add system prompt based on role
		systemPrompt := getTenantSystemPrompt()
		if role == "admin" {
			systemPrompt = getAdminSystemPrompt()
		}

		// Prepare messages with system prompt
		allMessages := make([]AIMessage, 0, len(messages)+1)
		allMessages = append(allMessages, AIMessage{
			Role:    "system",
			Content: systemPrompt,
		})
		allMessages = append(allMessages, messages...)

		requestBody := map[string]interface{}{
			"model":    ModelName,
			"messages": allMessages,
			"stream":   true,
		}

		jsonData, err := json.Marshal(requestBody)
		if err != nil {
			errorChan <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", BigModelBaseURL), bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

		client := &http.Client{Timeout: 120 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("failed to send request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("API error: %s", string(body))
			return
		}

		reader := resp.Body
		buf := make([]byte, 0, 1024)
		chunk := make([]byte, 1024)

		for {
			n, err := reader.Read(chunk)
			if n > 0 {
				buf = append(buf, chunk[:n]...)

				// Process complete lines
				for {
					lineEnd := -1
					for i := 0; i < len(buf); i++ {
						if buf[i] == '\n' {
							lineEnd = i
							break
						}
					}

					if lineEnd < 0 {
						break
					}

					line := string(buf[:lineEnd])
					buf = buf[lineEnd+1:]

					if len(line) > 6 && line[:6] == "data: " {
						data := line[6:]
						if data == "[DONE]" {
							return
						}

						var streamChunk AIStreamChunk
						if err := json.Unmarshal([]byte(data), &streamChunk); err != nil {
							continue
						}

						if len(streamChunk.Choices) > 0 && streamChunk.Choices[0].Delta.Content != "" {
							resultChan <- streamChunk.Choices[0].Delta.Content
						}
					}
				}
			}

			if err != nil {
				if err != io.EOF {
					errorChan <- fmt.Errorf("read error: %w", err)
				}
				break
			}
		}
	}()

	return resultChan, errorChan
}