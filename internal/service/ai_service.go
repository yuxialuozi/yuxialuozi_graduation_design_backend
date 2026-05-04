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

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const (
	BigModelBaseURL = "https://open.bigmodel.cn/api/paas/v4"
	ModelName       = "glm-4-flash"
)

type ChatResult struct {
	Content string
	Error   string
}

// Chat performs a non-streaming chat completion
func (s *AIService) Chat(role string, messages []AIMessage, context string) (*ChatResult, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("AI API key not configured, please set BIGMODEL_API_KEY environment variable")
	}

	// Build system prompt based on role
	systemPrompt := getSystemPrompt(role)
	if context != "" {
		systemPrompt += "\n\n【补充信息】\n" + context
	}

	allMessages := make([]AIMessage, 0, len(messages)+1)
	allMessages = append(allMessages, AIMessage{Role: "system", Content: systemPrompt})
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

func getSystemPrompt(role string) string {
	if role == "admin" {
		return getAdminSystemPrompt()
	}
	return getTenantSystemPrompt()
}

func getAdminSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的 AI 智能助手，专门帮助管理员进行数据分析和业务决策。

你的职责包括：
1. 数据分析：分析租户、合同、房间、费用、维修等数据，提供洞察和建议
2. 业务咨询：回答关于租金定价、合同条款、费用收取等业务问题
3. 趋势预测：基于历史数据，分析收入趋势、租户流失风险等
4. 维修建议：根据维修工单数据，提供维护优先级建议
5. 报表解读：帮助理解各类统计报表中的数据含义

【重要原则】：
1. 只回答与租户信息管理系统相关的业务问题
2. 不知道的问题明确告知用户，不编造答案
3. 涉及具体数据时，引导用户使用系统查看或说明如何查看
4. 回答应专业、准确、简洁

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04:05"))
}

func getTenantSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的 AI 智能助手，专门帮助租户（住户）解决日常问题。

你的职责包括：
1. 费用咨询：解答关于租金、水电费、物业费等费用问题
2. 合同解读：帮助理解租房合同的条款和权益
3. 维修指引：指导如何提交维修申请，跟踪维修进度
4. 政策说明：解释租金调整、押金退还、续租等政策
5. 常见问题：回答租户常见的各类问题

【重要原则】：
1. 只回答与租户信息管理系统相关的业务问题
2. 不知道的问题明确告知用户，不编造答案
3. 涉及具体个人数据时，引导用户登录系统查看
4. 回答应友好、耐心、通俗易懂
5. 不确定时，建议用户联系物业或在线客服

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04:05"))
}