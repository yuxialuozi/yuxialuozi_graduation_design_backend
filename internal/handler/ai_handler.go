package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"yuxialuozi_graduation_design_backend/internal/dto"
	"yuxialuozi_graduation_design_backend/internal/mcp"
	"yuxialuozi_graduation_design_backend/internal/service"
	"yuxialuozi_graduation_design_backend/pkg/response"
)

type AIHandler struct {
	aiService   *service.AIService
	toolExecutor *mcp.ToolExecutor
}

func NewAIHandler(aiService *service.AIService, toolExecutor *mcp.ToolExecutor) *AIHandler {
	return &AIHandler{
		aiService:   aiService,
		toolExecutor: toolExecutor,
	}
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req dto.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	userRole := "user"
	tenantID := uint(0)
	if r, exists := c.Get("role"); exists {
		if r == "admin" {
			userRole = "admin"
		}
	}
	if tid, exists := c.Get("tenantId"); exists {
		if t, ok := tid.(uint); ok {
			tenantID = t
		}
	}

	// Determine streaming or non-streaming
	if c.Query("stream") == "true" {
		h.streamChatWithMCP(c, userRole, tenantID, req.Messages)
		return
	}

	// Non-streaming (use existing service method without MCP)
	// For full MCP support, use streaming mode
	result, err := h.aiService.Chat(userRole, toServiceMessages(req.Messages), "")
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

// detectAndCallTools automatically detects intent and calls relevant tools
func (h *AIHandler) detectAndCallTools(userMessage string, tenantID uint) string {
	// Build a context string with all relevant data
	var context strings.Builder
	args := map[string]interface{}{"tenantId": float64(tenantID)}

	msg := strings.ToLower(userMessage)

	// Profile-related keywords (check early, before dashboard)
	profileKeywords := []string{"个人信息", "个人资料", "我的信息", "账号", "租户信息", "我的名字", "我的电话", "联系方式"}
	for _, kw := range profileKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_profile", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【个人信息】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	// Fee-related keywords (check early for specificity)
	feeKeywords := []string{"账单", "费用", "租金", "水电费", "物业费", "收费", "付款", "欠费", "未缴", "缴纳", "账单明细"}
	for _, kw := range feeKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_fees", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【费用数据】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	// Contract-related keywords
	contractKeywords := []string{"合同", "租期", "签约", "续租", "退租", "到期"}
	for _, kw := range contractKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_contracts", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【合同数据】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	// Maintenance-related keywords
	maintKeywords := []string{"维修", "报修", "工单", "修理", "损坏", "故障", "维护"}
	for _, kw := range maintKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_maintenance", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【维修工单数据】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	// Room-related keywords
	roomKeywords := []string{"房间", "房间号", "楼栋", "面积", "入住", "地址"}
	for _, kw := range roomKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_rooms", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【房间数据】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	// Dashboard-related keywords (check last - most generic)
	// Only match if no more specific keyword was detected
	dashKeywords := []string{"仪表盘", "概览", "统计", "汇总", "总结"}
	for _, kw := range dashKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("get_tenant_dashboard", args)
			if err == nil && result != "" && !strings.Contains(result, "error") {
				context.WriteString("\n【仪表盘数据】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			break
		}
	}

	return context.String()
}

// streamChatWithMCP handles streaming chat with MCP tool calling
func (h *AIHandler) streamChatWithMCP(c *gin.Context, role string, tenantID uint, messages []dto.AIMessage) {
	apiKey := os.Getenv("BIGMODEL_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI API key not configured"})
		return
	}

	// Auto-detect intent and call tools before sending to AI
	userMessage := ""
	for _, m := range messages {
		if m.Role == "user" {
			userMessage = m.Content
			break
		}
	}

	toolContext := h.detectAndCallTools(userMessage, tenantID)

	// Add tenant ID to args for tool calls
	args := map[string]interface{}{"tenantId": float64(tenantID)}

	// Build messages with system prompt
	allMessages := []service.AIMessage{
		{Role: "system", Content: getSystemPromptBaseWithContext(role, toolContext)},
	}
	allMessages = append(allMessages, toServiceMessages(messages)...)

	// First request with tools - use glm-4-flash (glm-4 has no credits)
	toolsJSON := mcp.GetToolsJSON(h.toolExecutor)
	requestBody := map[string]interface{}{
		"model":      "glm-4-flash",
		"messages":   allMessages,
		"stream":     true,
		"tools":      json.RawMessage(toolsJSON),
		"temperature": 0.3,
	}

	zap.L().Info("AI MCP request",
		zap.String("role", role),
		zap.Int("tenantId", int(tenantID)),
		zap.Int("toolsCount", 7),
		zap.Int("toolsJSONLen", len(toolsJSON)),
		zap.Int("messagesCount", len(allMessages)),
	)

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
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

	client := &http.Client{Timeout: 120 * time.Second}
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

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	// Collect all chunks and handle tool calls
	var collectedContent strings.Builder
	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				zap.L().Error("error reading stream", zap.Error(err))
			}
			break
		}
		c.Writer.WriteString(line)
		c.Writer.Flush()

		// Check if this is a tool call
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var chunk struct {
				Choices []struct {
					Delta struct {
						Content   string `json:"content"`
						ToolCalls []struct {
							ID       string `json:"id"`
							Type     string `json:"type"`
							Function struct {
								Name      string `json:"name"`
								Arguments string `json:"arguments"`
							} `json:"function"`
						} `json:"tool_calls"`
					} `json:"delta"`
					FinishReason string `json:"finish_reason"`
				} `json:"choices"`
			}

			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}

			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta
				if delta.Content != "" {
					collectedContent.WriteString(delta.Content)
				}

				// If AI wants to call a tool, handle it
				if len(delta.ToolCalls) > 0 && len(delta.ToolCalls[0].Function.Name) > 0 {
					// Flush any accumulated content before tool execution
					if collectedContent.Len() > 0 {
						c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", `{"choices":[{"delta":{"content":"`+fmt.Sprintf("%s\n\n[正在查询系统数据...]\n", collectedContent.String())+`"},"finish_reason":"stop"}]}`))
						c.Writer.Flush()
						collectedContent.Reset()
					}
					// Execute tool calls
					toolResults := h.executeToolCalls(delta.ToolCalls, args)
					// Add tool results as messages
					for _, tr := range toolResults {
						allMessages = append(allMessages, service.AIMessage{
							Role:    "tool",
							Content: fmt.Sprintf(`{"tool": "%s", "result": %s}`, tr.Name, tr.Result),
						})
					}
					// Continue conversation with tool results
					h.continueWithTools(c, client, apiKey, role, allMessages)
					return
				}
			}
		}
	}
}

// continueWithTools continues the conversation after tool execution
func (h *AIHandler) continueWithTools(c *gin.Context, client *http.Client, apiKey, role string, messages []service.AIMessage) {
	// Send continuation request with tool results
	requestBody := map[string]interface{}{
		"model":    "glm-4-flash",
		"messages": messages,
		"stream":   true,
		"tools":    json.RawMessage(mcp.GetToolsJSON(h.toolExecutor)),
	}

	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		c.Writer.WriteString(line)
		c.Writer.Flush()

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
		}
	}
}

// executeToolCalls executes multiple tool calls
func (h *AIHandler) executeToolCalls(toolCalls []struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}, args map[string]interface{}) []mcp.ToolResult {

	var results []mcp.ToolResult
	for _, tc := range toolCalls {
		// Parse arguments
		var toolArgs map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &toolArgs); err != nil {
			results = append(results, mcp.ToolResult{
				Name:  tc.Function.Name,
				Error: fmt.Sprintf("参数解析失败: %v", err),
			})
			continue
		}

		// Merge with tenant ID
		for k, v := range args {
			if _, exists := toolArgs[k]; !exists {
				toolArgs[k] = v
			}
		}

		// Execute tool
		result, err := h.toolExecutor.ExecuteTool(tc.Function.Name, toolArgs)
		if err != nil {
			results = append(results, mcp.ToolResult{
				Name:  tc.Function.Name,
				Error: err.Error(),
			})
		} else {
			results = append(results, mcp.ToolResult{
				Name:   tc.Function.Name,
				Result: result,
			})
		}
	}

	return results
}

func getSystemPromptBase(role string) string {
	if role == "admin" {
		return getAdminSystemPrompt()
	}
	return getTenantSystemPrompt()
}

func getSystemPromptBaseWithContext(role string, toolContext string) string {
	base := getSystemPromptBase(role)
	if toolContext == "" {
		return base
	}
	return base + "\n\n【系统实时数据】" + toolContext + "\n请基于以上真实数据回答用户问题。"
}

func getAdminSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的AI助手。你有7个工具可用：
1. get_tenant_dashboard - 获取仪表盘统计
2. get_tenant_fees - 获取费用账单
3. get_tenant_contracts - 获取合同列表
4. get_tenant_maintenance - 获取维修工单
5. get_tenant_rooms - 获取房间列表
6. search_knowledge - 搜索知识库
7. get_tenant_profile - 获取租户信息

【重要规则】当用户询问具体数据（如账单、合同、费用、维修、房间等）时，你必须先调用相关工具获取数据，然后基于真实数据回答。禁止凭空编造数据。

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04"))
}

func getTenantSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的AI助手。你有7个工具可用：
1. get_tenant_dashboard - 获取您的仪表盘统计
2. get_tenant_fees - 获取您的费用账单
3. get_tenant_contracts - 获取您的合同列表
4. get_tenant_maintenance - 获取您的维修工单
5. get_tenant_rooms - 获取您的房间信息
6. search_knowledge - 搜索知识库
7. get_tenant_profile - 获取您的个人信息

【重要规则】当用户询问具体数据（如账单、合同、费用、维修、房间等）时，你必须先调用相关工具获取数据，然后基于真实数据回答。禁止凭空编造数据。

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04"))
}

func toServiceMessages(messages []dto.AIMessage) []service.AIMessage {
	result := make([]service.AIMessage, len(messages))
	for i, msg := range messages {
		result[i] = service.AIMessage{Role: msg.Role, Content: msg.Content}
	}
	return result
}