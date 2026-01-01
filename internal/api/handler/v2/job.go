package v2

import (
	"citylife/internal/action"
	"citylife/internal/api/middleware"
	"citylife/internal/api/response"
	"citylife/internal/world"

	"github.com/gin-gonic/gin"
)

type jobInfo struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	DurationSeconds int            `json:"duration_seconds"`
	Income          int            `json:"income"`
	Consume         map[string]int `json:"consume,omitempty"`
}

// GetJobs 获取工作列表（结构化）
func GetJobs(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	actions := action.GetJobActions()
	jobs := make([]jobInfo, 0, len(actions))
	for _, act := range actions {
		work, ok := act.(*action.WorkAction)
		if !ok {
			continue
		}
		jobs = append(jobs, jobInfo{
			ID:              work.ID(),
			Name:            work.Name(),
			DurationSeconds: work.DurationSeconds(),
			Income:          work.Income(),
			Consume:         work.Consume(),
		})
	}

	ok(c, jobs)
}

// DoJob 执行指定工作（需要在人才市场）
func DoJob(c *gin.Context) {
	sess := middleware.GetSession(c)
	if sess == nil {
		response.SessionNotFound(c)
		return
	}

	state := sess.State
	if !state.User.IsAlive() {
		failGameOver(c, state)
		return
	}

	if state.World.Where != world.LocationJobMarket {
		fail(c, response.ErrCodeInvalidLocation, "你需要在人才市场才能找工作", actionData{State: buildState(state)})
		return
	}

	jobID := c.Param("id")
	var target action.Action
	for _, act := range action.GetJobActions() {
		if act.ID() == jobID {
			target = act
			break
		}
	}

	if target == nil {
		response.BadRequest(c, "无效的工作ID")
		return
	}

	executeAndRespond(c, state, target, "")
}
