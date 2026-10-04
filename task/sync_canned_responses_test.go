package task

import (
	"testing"

	"github.com/twbworld/agent/global"
	"github.com/twbworld/agent/model/config"
	"github.com/twbworld/agent/model/enum"
)

func TestParseShortCode(t *testing.T) {
	// 模拟全局配置
	global.Config = &config.Config{
		Ai: config.Ai{
			SemanticPrefix: "ai@",
			HybridPrefix:   "ai+@",
		},
	}

	m := &Manager{}

	tests := []struct {
		name         string
		shortCode    string
		expectedType enum.KeywordType
		expectedText string
	}{
		{
			name:         "纯精确匹配",
			shortCode:    "人工服务",
			expectedType: enum.KeywordTypeExact,
			expectedText: "人工服务",
		},
		{
			name:         "纯语义匹配",
			shortCode:    "ai@退款怎么退",
			expectedType: enum.KeywordTypeSemantic,
			expectedText: "退款怎么退",
		},
		{
			name:         "混合匹配(精确+语义)",
			shortCode:    "ai+@发货时间",
			expectedType: enum.KeywordTypeHybrid,
			expectedText: "发货时间",
		},
		{
			name:         "包含前缀字符但并非规定前缀的普通匹配",
			shortCode:    "ai你好",
			expectedType: enum.KeywordTypeExact,
			expectedText: "ai你好",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qType, qText := m.parseShortCode(tt.shortCode)
			if qType != tt.expectedType || qText != tt.expectedText {
				t.Errorf("parseShortCode() 失败!\n输入: %s\n预期: [%s] %s\n实际: [%s] %s",
					tt.shortCode, tt.expectedType, tt.expectedText, qType, qText)
			}
		})
	}
}
