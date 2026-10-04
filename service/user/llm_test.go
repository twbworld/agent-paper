package user

import (
	"testing"
)

func TestCleanJSON(t *testing.T) {
	svc := &llmService{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "标准JSON没有多余内容",
			input:    `{"intent": "order_inquiry", "emotion": "neutral"}`,
			expected: `{"intent": "order_inquiry", "emotion": "neutral"}`,
		},
		{
			name:     "带有Markdown包裹的JSON",
			input:    "```json\n{\"intent\": \"order_inquiry\", \"emotion\": \"neutral\"}\n```",
			expected: `{"intent": "order_inquiry", "emotion": "neutral"}`,
		},
		{
			name:     "仅带有反引号的JSON",
			input:    "```\n{\"intent\": \"product_inquiry\"}\n```",
			expected: `{"intent": "product_inquiry"}`,
		},
		{
			name:     "带有前后空格和换行",
			input:    "  \n\n  ```json\n{\"urgency\": \"high\"}\n``` \n ",
			expected: `{"urgency": "high"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := svc.cleanJSON(tt.input)
			if actual != tt.expected {
				t.Errorf("cleanJSON() 失败!\n预期: %s\n实际: %s", tt.expected, actual)
			}
		})
	}
}
