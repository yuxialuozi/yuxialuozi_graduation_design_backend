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
	aiService    *service.AIService
	toolExecutor *mcp.ToolExecutor
}

func NewAIHandler(aiService *service.AIService, toolExecutor *mcp.ToolExecutor) *AIHandler {
	return &AIHandler{
		aiService:    aiService,
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

	// Non-streaming
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
func (h *AIHandler) detectAndCallTools(userMessage string, role string, tenantID uint) string {
	var context strings.Builder
	msg := strings.ToLower(userMessage)

	if role == "admin" {
		adminKeywords := map[string][]string{
			"get_admin_dashboard":   {"仪表盘", "概览", "统计", "汇总", "总结", "入住率", "总租户"},
			"get_admin_fees":        {"账单", "费用", "收费", "欠费", "未缴", "收入", "缴费"},
			"get_admin_tenants":     {"租户", "租户列表", "租户排名", "租户统计"},
			"get_admin_contracts":   {"合同", "租期", "签约", "到期", "合同统计"},
			"get_admin_maintenance": {"维修", "报修", "工单", "修理", "故障", "维护"},
		}

		for toolName, keywords := range adminKeywords {
			for _, kw := range keywords {
				if strings.Contains(msg, kw) {
					zap.L().Info("admin keyword matched", zap.String("tool", toolName), zap.String("keyword", kw))
					result, err := h.toolExecutor.ExecuteTool(toolName, map[string]interface{}{})
					zap.L().Info("tool execution result", zap.String("tool", toolName), zap.Error(err), zap.Int("resultLen", len(result)))
					if err == nil && result != "" {
						context.WriteString("\n【系统数据】\n")
						context.WriteString(result)
						context.WriteString("\n")
					}
					return context.String()
				}
			}
		}
	} else {
		args := map[string]interface{}{"tenantId": float64(tenantID)}

		tenantKeywords := map[string][]string{
			"get_tenant_profile":     {"个人信息", "个人资料", "我的信息", "账号", "租户信息", "我的名字", "我的电话", "联系方式"},
			"get_tenant_fees":        {"账单", "费用", "租金", "水电费", "物业费", "收费", "付款", "欠费", "未缴", "缴纳", "账单明细"},
			"get_tenant_contracts":   {"合同", "租期", "签约", "续租", "退租", "到期"},
			"get_tenant_maintenance": {"维修", "报修", "工单", "修理", "损坏", "故障", "维护"},
			"get_tenant_rooms":       {"房间", "房间号", "楼栋", "面积", "入住", "地址"},
			"get_tenant_dashboard":   {"仪表盘", "概览", "统计", "汇总", "总结"},
		}

		for toolName, keywords := range tenantKeywords {
			for _, kw := range keywords {
				if strings.Contains(msg, kw) {
					result, err := h.toolExecutor.ExecuteTool(toolName, args)
					if err == nil && result != "" {
						label := map[string]string{
							"get_tenant_profile":     "个人信息",
							"get_tenant_fees":        "费用数据",
							"get_tenant_contracts":   "合同数据",
							"get_tenant_maintenance": "维修工单数据",
							"get_tenant_rooms":       "房间数据",
							"get_tenant_dashboard":   "仪表盘数据",
						}
						context.WriteString("\n【" + label[toolName] + "】\n")
						context.WriteString(result)
						context.WriteString("\n")
					}
					return context.String()
				}
			}
		}
	}

	// Fallback: search knowledge base
	knowledgeKeywords := []string{"怎么", "如何", "流程", "政策", "规定", "说明", "什么"}
	for _, kw := range knowledgeKeywords {
		if strings.Contains(msg, kw) {
			result, err := h.toolExecutor.ExecuteTool("search_knowledge", map[string]interface{}{"query": userMessage})
			if err == nil && result != "" {
				context.WriteString("\n【知识库】\n")
				context.WriteString(result)
				context.WriteString("\n")
			}
			return context.String()
		}
	}

	return context.String()
}

// toolCallInfo stores tool call data from AI response
type toolCallInfo struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
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

	zap.L().Info("streamChatWithMCP start", zap.String("role", role), zap.Uint("tenantId", tenantID), zap.String("userMessage", userMessage))
	toolContext := h.detectAndCallTools(userMessage, role, tenantID)
	zap.L().Info("intent detection result", zap.String("role", role), zap.String("toolContext", toolContext), zap.Int("toolContextLen", len(toolContext)))

	// Build messages with system prompt - use data-aware prompt when tool context is available
	var systemPrompt string
	if toolContext != "" {
		systemPrompt = getSystemPromptWithData(role, toolContext)
	} else {
		systemPrompt = getSystemPromptBaseWithContext(role, toolContext)
	}
	allMessages := []map[string]interface{}{
		{"role": "system", "content": systemPrompt},
	}
	for _, m := range messages {
		allMessages = append(allMessages, map[string]interface{}{
			"role":    m.Role,
			"content": m.Content,
		})
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)

	// If we have tool context from intent detection, send directly without tools
	// The data is already injected into the system prompt
	if toolContext != "" {
		h.sendToBigModelDirect(c, apiKey, allMessages)
		return
	}

	// No tool context detected - send with tools for AI to decide
	toolsJSON := mcp.GetToolsJSONByRole(h.toolExecutor, role)
	h.sendToBigModelWithTools(c, apiKey, role, tenantID, allMessages, toolsJSON)
}

// sendToBigModelDirect sends a direct request (no tools) to BigModel
func (h *AIHandler) sendToBigModelDirect(c *gin.Context, apiKey string, messages []map[string]interface{}) {
	requestBody := map[string]interface{}{
		"model":       "glm-4-flash",
		"messages":    messages,
		"stream":      true,
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		h.sendSSEError(c, "请求创建失败")
		return
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		h.sendSSEError(c, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to send direct request to BigModel", zap.Error(err))
		h.sendSSEError(c, "AI 服务连接失败")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		zap.L().Error("BigModel direct API error", zap.String("body", string(bodyBytes)))
		h.sendSSEError(c, "AI 服务错误")
		return
	}

	// Stream response directly to client
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		c.Writer.WriteString(line)
		c.Writer.Flush()
	}
}

// sendToBigModelWithTools sends a request to BigModel and handles tool calls
func (h *AIHandler) sendToBigModelWithTools(c *gin.Context, apiKey string, role string, tenantID uint, messages []map[string]interface{}, toolsJSON string) {
	requestBody := map[string]interface{}{
		"model":       "glm-4-flash",
		"messages":    messages,
		"stream":      true,
		"tools":       json.RawMessage(toolsJSON),
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		h.sendSSEError(c, "failed to create request")
		return
	}

	req, err := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	if err != nil {
		h.sendSSEError(c, err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to send request to BigModel", zap.Error(err))
		h.sendSSEError(c, "failed to connect to AI service")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		zap.L().Error("BigModel API error", zap.Int("status", resp.StatusCode), zap.String("body", string(bodyBytes)))
		h.sendSSEError(c, "AI service error")
		return
	}

	// Read the full response, collecting tool calls and content
	var contentBuilder strings.Builder
	var toolCalls []toolCallInfo
	var hasToolCalls bool

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				zap.L().Error("error reading stream", zap.Error(err))
			}
			break
		}

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   string         `json:"content"`
					ToolCalls []toolCallInfo `json:"tool_calls"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		choice := chunk.Choices[0]
		delta := choice.Delta

		// Collect content
		if delta.Content != "" {
			contentBuilder.WriteString(delta.Content)
		}

		// Collect tool calls
		for _, tc := range delta.ToolCalls {
			if tc.Function.Name != "" {
				toolCalls = append(toolCalls, tc)
				hasToolCalls = true
			}
		}

		// When finish_reason indicates tool calls or stream is done
		if choice.FinishReason == "tool_calls" || choice.FinishReason == "stop" {
			// If tool calls were collected, execute them and continue
			if hasToolCalls && len(toolCalls) > 0 {
				zap.L().Info("executing tool calls", zap.Int("count", len(toolCalls)))
				h.handleToolCallsAndContinue(c, client, apiKey, role, tenantID, messages, contentBuilder.String(), toolCalls, toolsJSON)
				return
			}
		}
	}

	// No tool calls - just stream the collected content to frontend
	content := contentBuilder.String()
	if content != "" {
		h.sendSSEContent(c, content)
	}
	h.sendSSEDone(c)
}

// handleToolCallsAndContinue executes tools and continues the conversation
func (h *AIHandler) handleToolCallsAndContinue(c *gin.Context, client *http.Client, apiKey string, role string, tenantID uint, messages []map[string]interface{}, preContent string, toolCalls []toolCallInfo, toolsJSON string) {
	// Execute tool calls
	toolResults := h.executeToolCalls(toolCalls, role, tenantID)

	// Build assistant message with tool_calls for BigModel API format
	apiToolCalls := make([]map[string]interface{}, len(toolCalls))
	for i, tc := range toolCalls {
		apiToolCalls[i] = map[string]interface{}{
			"id":   tc.ID,
			"type": "function",
			"function": map[string]interface{}{
				"name":      tc.Function.Name,
				"arguments": tc.Function.Arguments,
			},
		}
	}

	// Add assistant message with tool_calls (content must be null for tool calls)
	messages = append(messages, map[string]interface{}{
		"role":       "assistant",
		"content":    nil,
		"tool_calls": apiToolCalls,
	})

	// Add tool result messages
	for _, tr := range toolResults {
		toolCallID := ""
		for _, tc := range toolCalls {
			if tc.Function.Name == tr.Name {
				toolCallID = tc.ID
				break
			}
		}
		resultContent := tr.Result
		if tr.Error != "" {
			resultContent = fmt.Sprintf(`{"error": "%s"}`, tr.Error)
		}
		messages = append(messages, map[string]interface{}{
			"role":         "tool",
			"content":      resultContent,
			"tool_call_id": toolCallID,
		})
	}

	// Send "querying data" indicator then continue
	zap.L().Info("tool results messages", zap.Int("count", len(messages)), zap.String("lastRole", fmt.Sprintf("%v", messages[len(messages)-1]["role"])))

	if preContent != "" {
		h.sendSSEContent(c, preContent+"\n\n[正在查询系统数据...]\n\n")
	} else {
		h.sendSSEContent(c, "[正在查询系统数据...]\n\n")
	}

	// Continue the conversation with tool results (no more tools to avoid infinite loops)
	h.sendToBigModelFinal(c, apiKey, role, messages)
}

// sendToBigModelFinal sends a final request (without tools) to get the natural language response
func (h *AIHandler) sendToBigModelFinal(c *gin.Context, apiKey string, role string, messages []map[string]interface{}) {
	requestBody := map[string]interface{}{
		"model":       "glm-4-flash",
		"messages":    messages,
		"stream":      true,
		"temperature": 0.3,
	}

	jsonBody, _ := json.Marshal(requestBody)
	zap.L().Info("final request body", zap.Int("len", len(jsonBody)), zap.String("body", string(jsonBody)[:min(len(jsonBody), 500)]))
	req, _ := http.NewRequest("POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to send final request to BigModel", zap.Error(err))
		h.sendSSEError(c, "AI 服务连接失败")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		zap.L().Error("BigModel final API error", zap.String("body", string(bodyBytes)))
		h.sendSSEError(c, "AI 服务错误")
		return
	}

	// Stream the final response directly to the client
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		c.Writer.WriteString(line)
		c.Writer.Flush()
	}
}

// sendSSEContent sends a content chunk to the client via SSE
func (h *AIHandler) sendSSEContent(c *gin.Context, content string) {
	escaped, _ := json.Marshal(content)
	sseData := fmt.Sprintf(`{"choices":[{"delta":{"content":%s},"finish_reason":null}]}`, string(escaped))
	fmt.Fprintf(c.Writer, "data: %s\n\n", sseData)
	c.Writer.Flush()
}

// sendSSEError sends an error message to the client via SSE
func (h *AIHandler) sendSSEError(c *gin.Context, errMsg string) {
	escaped, _ := json.Marshal("⚠️ " + errMsg)
	sseData := fmt.Sprintf(`{"choices":[{"delta":{"content":%s},"finish_reason":"stop"}]}`, string(escaped))
	fmt.Fprintf(c.Writer, "data: %s\n\n", sseData)
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
}

// sendSSEDone sends the [DONE] marker
func (h *AIHandler) sendSSEDone(c *gin.Context) {
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	c.Writer.Flush()
}

// executeToolCalls executes multiple tool calls
func (h *AIHandler) executeToolCalls(toolCalls []toolCallInfo, role string, tenantID uint) []mcp.ToolResult {
	var results []mcp.ToolResult

	// Build base args based on role
	baseArgs := map[string]interface{}{}
	if role == "user" {
		baseArgs["tenantId"] = float64(tenantID)
	}

	for _, tc := range toolCalls {
		var toolArgs map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &toolArgs); err != nil {
			toolArgs = map[string]interface{}{}
		}

		// Merge base args
		for k, v := range baseArgs {
			if _, exists := toolArgs[k]; !exists {
				toolArgs[k] = v
			}
		}

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

func getSystemPromptWithData(role string, toolContext string) string {
	if role == "admin" {
		return fmt.Sprintf(`你是租户信息管理系统的AI智能助手，专门帮助管理员进行数据分析和业务决策。

【系统实时数据】%s

请基于以上真实数据回答用户问题。回答要求：
1. 基于真实数据，禁止编造任何数据
2. 回答应专业、准确，提供数据洞察和建议
3. 涉及趋势分析时，基于真实数据给出合理判断

当前系统时间：%s`, toolContext, time.Now().Format("2006-01-02 15:04"))
	}
	return fmt.Sprintf(`你是租户信息管理系统的AI智能助手，专门帮助租户（住户）解决日常问题。

【系统实时数据】%s

请基于以上真实数据回答用户问题。回答要求：
1. 基于真实数据，禁止编造任何数据
2. 回答应友好、耐心、通俗易懂
3. 不确定时，建议用户联系物业或在线客服

当前系统时间：%s`, toolContext, time.Now().Format("2006-01-02 15:04"))
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
	return fmt.Sprintf(`你是租户信息管理系统的AI智能助手，专门帮助管理员进行数据分析和业务决策。

你可以调用以下工具获取系统数据：
1. get_admin_dashboard - 获取系统全局仪表盘数据（入住率、收入、合同、维修统计等）
2. get_admin_fees - 获取系统全局费用统计和明细
3. get_admin_contracts - 获取合同统计和明细
4. get_admin_tenants - 获取租户统计和排名
5. get_admin_maintenance - 获取维修工单统计和明细
6. search_knowledge - 搜索知识库

【重要规则】
1. 当用户询问具体数据时，你必须先调用相关工具获取真实数据，然后基于真实数据回答。禁止凭空编造数据。
2. 作为管理员助手，你的回答应专业、准确，提供数据洞察和建议。
3. 涉及趋势分析时，基于真实数据给出合理判断。

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04"))
}

func getTenantSystemPrompt() string {
	return fmt.Sprintf(`你是租户信息管理系统的AI智能助手，专门帮助租户（住户）解决日常问题。

你可以调用以下工具获取个人数据：
1. get_tenant_profile - 获取您的个人信息
2. get_tenant_fees - 获取您的费用账单
3. get_tenant_contracts - 获取您的合同列表
4. get_tenant_maintenance - 获取您的维修工单
5. get_tenant_rooms - 获取您的房间信息
6. get_tenant_dashboard - 获取您的仪表盘统计
7. search_knowledge - 搜索知识库

【重要规则】
1. 当用户询问具体数据时，你必须先调用相关工具获取真实数据，然后基于真实数据回答。禁止凭空编造数据。
2. 作为租户助手，你的回答应友好、耐心、通俗易懂。
3. 不确定时，建议用户联系物业或在线客服。

当前系统时间：%s`, time.Now().Format("2006-01-02 15:04"))
}

func toServiceMessages(messages []dto.AIMessage) []service.AIMessage {
	result := make([]service.AIMessage, len(messages))
	for i, msg := range messages {
		result[i] = service.AIMessage{Role: msg.Role, Content: msg.Content}
	}
	return result
}
