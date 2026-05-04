package handler

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type AIHandler struct {
	aiService *service.AIService
}

func NewAIHandler(aiService *service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// Chat godoc
// @Summary AI 聊天
// @Description 与 AI 助手对话，支持流式和非流式响应
// @Tags AI助手
// @Accept json
// @Produce json, text/event-stream
// @Security BearerAuth
// @Param request body dto.AIChatRequest true "聊天请求"
// @Param stream query bool false "是否流式响应" default(false)
// @Success 200 {object} dto.AIChatResponse "非流式响应"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	var req dto.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	// Determine user role from JWT claims
	userRole := "user"
	if r, exists := c.Get("role"); exists {
		if r == "admin" {
			userRole = "admin"
		}
	}

	// Check if streaming is requested
	if c.Query("stream") == "true" {
		h.streamChat(c, userRole, req.Messages)
		return
	}

	// Non-streaming response
	result, err := h.aiService.Chat(userRole, toServiceMessages(req.Messages))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	if result.Error != "" {
		response.InternalError(c, result.Error)
		return
	}

	response.Success(c, dto.AIChatResponse{Content: result.Content})
}

// streamChat handles streaming chat responses
func (h *AIHandler) streamChat(c *gin.Context, role string, messages []dto.AIMessage) {
	apiKey := os.Getenv("BIGMODEL_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI API key not configured, please set BIGMODEL_API_KEY environment variable"})
		return
	}

	// Prepare messages with system prompt
	allMessages := []service.AIMessage{}
	systemPrompt := getTenantSystemPrompt()
	if role == "admin" {
		systemPrompt = getAdminSystemPrompt()
	}
	allMessages = append(allMessages, service.AIMessage{Role: "system", Content: systemPrompt})
	allMessages = append(allMessages, toServiceMessages(messages)...)

	// Create request body
	requestBody := map[string]interface{}{
		"model":    "glm-4-flash",
		"messages": allMessages,
		"stream":   true,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		zap.L().Error("failed to marshal request body", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to send request to BigModel", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to connect to AI service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service error: " + string(bodyBytes)})
		return
	}

	// Set SSE headers for streaming
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Status(http.StatusOK)

	// Stream response to client
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				zap.L().Error("error reading stream", zap.Error(err))
			}
			break
		}

		// Forward SSE data to client
		c.Writer.WriteString(line)
		c.Writer.Flush()
	}
}

// getAdminSystemPrompt returns the system prompt for admin users
func getAdminSystemPrompt() string {
	return `你是租户信息管理系统的 AI 智能助手，专门帮助管理员进行数据分析和业务决策。

你的职责包括：
1. 数据分析：分析租户、合同、房间、费用、维修等数据，提供洞察和建议
2. 业务咨询：回答关于租金定价、合同条款、费用收取等业务问题
3. 趋势预测：基于历史数据，分析收入趋势、租户流失风险等
4. 维修建议：根据维修工单数据，提供维护优先级建议
5. 报表解读：帮助理解各类统计报表中的数据含义

请用专业、简洁的语言回答，如果有具体数据可以给出量化分析。`
}

// getTenantSystemPrompt returns the system prompt for tenant users
func getTenantSystemPrompt() string {
	return `你是租户信息管理系统的 AI 智能助手，专门帮助租户（住户）解决日常问题。

你的职责包括：
1. 费用咨询：解答关于租金、水电费、物业费等费用问题
2. 合同解读：帮助理解租房合同的条款和权益
3. 维修指引：指导如何提交维修申请，跟踪维修进度
4. 政策说明：解释租金调整、押金退还、续租等政策
5. 常见问题：回答租户常见的各类问题

请用友好、耐心的语言回答，使用通俗易懂的语言解释专业术语。`
}

// toServiceMessages converts DTO messages to service messages
func toServiceMessages(messages []dto.AIMessage) []service.AIMessage {
	result := make([]service.AIMessage, len(messages))
	for i, msg := range messages {
		result[i] = service.AIMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return result
}