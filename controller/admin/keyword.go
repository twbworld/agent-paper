package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/twbworld/agent/model/common"
	"github.com/twbworld/agent/service"
)

type KeywordApi struct{}

func (k *KeywordApi) ListItems(c *gin.Context) {
	items, err := service.Service.AdminServiceGroup.KeywordService.ListItems(c)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}
	common.Success(c, items)
}

func (k *KeywordApi) UpsertItem(c *gin.Context) {
	var req common.UpsertKnowledgeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, err.Error())
		return
	}

	if err := service.Service.AdminServiceGroup.KeywordService.UpsertItem(c, &req); err != nil {
		common.Fail(c, err.Error())
		return
	}
	common.Success(c, nil)
}

func (k *KeywordApi) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")
	if itemID == "" {
		common.Fail(c, "ID 不能为空")
		return
	}

	if err := service.Service.AdminServiceGroup.KeywordService.DeleteItem(c, itemID); err != nil {
		common.Fail(c, err.Error())
		return
	}
	common.Success(c, nil)
}

func (k *KeywordApi) GenerateQuestions(c *gin.Context) {
	var req common.GenerateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, err.Error())
		return
	}

	resp, err := service.Service.AdminServiceGroup.KeywordService.GenerateQuestions(c, &req)
	if err != nil {
		common.Fail(c, err.Error())
		return
	}
	common.Success(c, resp)
}
