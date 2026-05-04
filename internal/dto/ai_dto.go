package dto

// AIChatRequest represents a chat message from frontend
type AIChatRequest struct {
	Messages []AIMessage `json:"messages" binding:"required"`
}

// AIMessage represents a chat message
type AIMessage struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// AIChatResponse represents the chat response
type AIChatResponse struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}