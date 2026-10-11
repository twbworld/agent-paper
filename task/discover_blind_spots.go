package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sashabaranov/go-openai"
	"github.com/twbworld/agent/global"
	"github.com/twbworld/agent/internal/chatwoot"
	"github.com/twbworld/agent/internal/llm"
	"github.com/twbworld/agent/model/common"
	"github.com/twbworld/agent/model/enum"
)

// DiscoverBlindSpots 从转人工的会话记录中提炼未解决的问题，并以占位符形式写回知识库
func (m *Manager) DiscoverBlindSpots() error {
	global.Log.Info("开始执行闲时知识盲区自动挖掘任务...")

	ctx := context.Background()
	if global.RedisClient == nil || global.ChatwootService == nil || global.LlmService == nil {
		return fmt.Errorf("依赖服务未初始化，任务中止")
	}

	// 1. 获取所有待复盘的转人工会话队列
	queue, err := global.RedisClient.HGetAll(ctx, "kb:blind_spots_queue").Result()
	if err != nil || len(queue) == 0 {
		global.Log.Info("当前无待挖掘的知识盲区队列")
		return nil
	}

	accountID := uint(global.Config.Chatwoot.AccountId)
	semanticPrefix := global.Config.Ai.SemanticPrefix

	for convIDStr := range queue {
		convID, _ := strconv.ParseUint(convIDStr, 10, 32)
		if convID == 0 {
			continue
		}

		// 2. 调用 Chatwoot API 拉取真实的完整历史记录
		msgs, err := global.ChatwootService.GetConversationMessages(ctx, accountID, uint(convID))
		if err != nil {
			global.Log.Warnf("拉取会话 %d 历史失败: %v", convID, err)
			continue
		}

		// 组装用于 LLM 分析的纯文本对话日志
		var chatLog strings.Builder
		for i := len(msgs) - 1; i >= 0; i-- { // 确保按时间顺序拼接
			msg := msgs[i]
			if msg.Content == "" {
				continue
			}
			role := "用户"
			if msg.MessageType == chatwoot.MessageDirectionOutgoing {
				role = "客服/AI"
			}
			chatLog.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
		}

		if chatLog.Len() == 0 {
			global.RedisClient.HDel(ctx, "kb:blind_spots_queue", convIDStr)
			continue
		}

		// 3. 调用小模型提炼未解决的核心痛点
		prompt := "你是一个知识库提炼助手。阅读以下最终导致转人工的客服对话记录。如果用户有明确的业务痛点但由于AI知识库匮乏未能解决，请用一句话提炼出这个核心问题（限20字内）。如果没有明确的新业务问题或仅是情绪发泄，请严格回复'无'。"

		respMsg, err := global.LlmService.Chat(ctx, &llm.ChatRequest{
			Size:         enum.ModelSmall,
			SystemPrompt: enum.SystemPrompt(prompt),
			Messages: []common.LlmMessage{
				{Role: openai.ChatMessageRoleUser, Content: chatLog.String()},
			},
			Temperature:     new(float32(0.1)), // 低温保证提炼的确定性
			DisableThinking: true,
		})

		if err != nil {
			global.Log.Warnf("挖掘会话 %d 盲区失败: %v", convID, err)
			continue
		}

		extractedQuestion := strings.Trim(strings.TrimSpace(respMsg.Content), `"'。， `)

		// 4. 验证提炼结果并回写知识库
		if extractedQuestion != "" && extractedQuestion != "无" && !strings.Contains(extractedQuestion, "转人工") {

			shortCode := semanticPrefix + extractedQuestion
			// 巧妙利用特殊标识作为 Answer 入库，前端直接展示，管理员一目了然
			placeholderAnswer := "【知识盲区待完善】这是系统夜间自动挖掘到的未解决问题。请管理员点击编辑，在此处补充真实答案，随后点击“AI生成提问”并保存即可生效。"

			_, err := global.ChatwootService.CreateCannedResponse(ctx, shortCode, placeholderAnswer)
			if err == nil {
				global.Log.Infof("成功挖掘并自动注入知识盲区待办项: %s", extractedQuestion)
			} else {
				global.Log.Warnf("注入知识盲区失败: %v", err)
			}
		}

		// 无论是否提炼出问题，处理完毕后都将其移出队列
		global.RedisClient.HDel(ctx, "kb:blind_spots_queue", convIDStr)
	}

	global.Log.Info("闲时知识盲区自动挖掘任务执行完毕")
	return nil
}
