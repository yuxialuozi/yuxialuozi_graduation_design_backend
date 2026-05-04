package mcp

import (
	"encoding/json"
)

// Tool represents an MCP tool that the AI can call
type Tool struct {
	Name        string
	Description string
	Parameters  ToolParameters
	Handler     ToolHandler
}

// ToolParameters defines the parameters for a tool
type ToolParameters struct {
	Type       string
	Properties map[string]ToolProperty
	Required   []string
}

// ToolProperty defines a single parameter property
type ToolProperty struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// ToolHandler is the function that executes a tool
type ToolHandler func(args map[string]interface{}) (string, error)

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Name   string
	Result string
	Error  string
}

// ToFunctionDef converts a tool to BigModel API function definition format
func (t *Tool) ToFunctionDef() map[string]interface{} {
	return map[string]interface{}{
		"type": "function",
		"function": map[string]interface{}{
			"name":        t.Name,
			"description": t.Description,
			"parameters": map[string]interface{}{
				"type":       "object",
				"properties": t.Parameters.Properties,
				"required":   t.Parameters.Required,
			},
		},
	}
}

// GetAllTools returns all available MCP tools
func GetAllTools(executor *ToolExecutor) []Tool {
	return []Tool{
		{
			Name:        "get_tenant_profile",
			Description: "获取当前登录租户的个人信息，包括姓名、联系方式、入住房间数、有效合同数等",
			Parameters: ToolParameters{
				Type:       "object",
				Properties: map[string]ToolProperty{},
				Required:   []string{},
			},
			Handler: executor.GetTenantProfile,
		},
		{
			Name:        "get_tenant_fees",
			Description: "获取当前租户的费用账单列表，包括租金、水电费、物业费等，支持筛选未缴/已缴状态",
			Parameters: ToolParameters{
				Type: "object",
				Properties: map[string]ToolProperty{
					"status": {
						Type:        "string",
						Description: "账单状态筛选：paid=已缴, unpaid=未缴, overdue=逾期，不传则返回全部",
						Enum:        []string{"paid", "unpaid", "overdue"},
					},
					"page": {
						Type:        "integer",
						Description: "页码，默认1",
					},
					"pageSize": {
						Type:        "integer",
						Description: "每页数量，默认10",
					},
				},
				Required: []string{},
			},
			Handler: executor.GetTenantFees,
		},
		{
			Name:        "get_tenant_contracts",
			Description: "获取当前租户的所有合同列表，包括合同编号、租期、金额、状态等",
			Parameters: ToolParameters{
				Type: "object",
				Properties: map[string]ToolProperty{
					"status": {
						Type:        "string",
						Description: "合同状态筛选：draft=草稿, active=生效中, expired=已过期, terminated=已终止",
						Enum:        []string{"draft", "active", "expired", "terminated"},
					},
				},
				Required: []string{},
			},
			Handler: executor.GetTenantContracts,
		},
		{
			Name:        "get_tenant_maintenance",
			Description: "获取当前租户的维修工单列表，包括工单号、类型、状态、处理进度等",
			Parameters: ToolParameters{
				Type: "object",
				Properties: map[string]ToolProperty{
					"status": {
						Type:        "string",
						Description: "工单状态筛选：pending=待处理, processing=处理中, completed=已完成, cancelled=已取消",
						Enum:        []string{"pending", "processing", "completed", "cancelled"},
					},
				},
				Required: []string{},
			},
			Handler: executor.GetTenantMaintenance,
		},
		{
			Name:        "get_tenant_dashboard",
			Description: "获取当前租户的仪表盘数据，包括有效合同数、租住房间数、累计费用、未缴费用、待处理维修数等统计信息",
			Parameters: ToolParameters{
				Type:       "object",
				Properties: map[string]ToolProperty{},
				Required:   []string{},
			},
			Handler: executor.GetTenantDashboard,
		},
		{
			Name:        "get_tenant_rooms",
			Description: "获取当前租户入住的房间列表，包括房间号、楼栋、面积、月租等",
			Parameters: ToolParameters{
				Type:       "object",
				Properties: map[string]ToolProperty{},
				Required:   []string{},
			},
			Handler: executor.GetTenantRooms,
		},
		{
			Name:        "search_knowledge",
			Description: "搜索租户管理系统的知识库，查找相关的业务知识和政策说明",
			Parameters: ToolParameters{
				Type: "object",
				Properties: map[string]ToolProperty{
					"query": {
						Type:        "string",
						Description: "搜索关键词，如：租金计算、水电费、维修流程、合同条款等",
					},
				},
				Required: []string{"query"},
			},
			Handler: executor.SearchKnowledge,
		},
	}
}

// GetToolsJSON returns all tools as JSON for BigModel API
func GetToolsJSON(executor *ToolExecutor) string {
	tools := GetAllTools(executor)
	result := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		result[i] = tool.ToFunctionDef()
	}
	data, _ := json.Marshal(result)
	return string(data)
}