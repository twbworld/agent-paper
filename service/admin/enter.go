package admin

import "github.com/twbworld/agent/task"

type ServiceGroup struct {
	KeywordService   KeywordService
	DashboardService DashboardService
}

func NewServiceGroup(taskManager *task.Manager) ServiceGroup {
	return ServiceGroup{
		KeywordService:   NewKeywordService(taskManager),
		DashboardService: NewDashboardService(),
	}
}
