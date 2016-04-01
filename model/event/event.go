package event

import (
	"fmt"
	"github.com/Cepave/alarm/logger"
	coommonModel "github.com/Cepave/common/model"
	"github.com/Cepave/common/utils"
	"github.com/astaxie/beego/orm"
	"time"
)

type EventCases struct {
	// uniuq
	Id       string `json:"id" orm:"pk"`
	Endpoint string `json:"endpoint"`
	Metric   string `json:"metric"`
	Func     string `json:"func"`
	Cond     string `json:"cond"`
	Note     string `json:"note"`
	//leftValue + operator + rightValue
	MaxStep      int       `json:"max_step"`
	CurrentStep  int       `json:"current_step"`
	Priority     int       `json:"priority"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"start_at"`
	UpdateAt     time.Time `json:"update_at"`
	ClosedAt     time.Time `json:"closed_at"`
	ClosedNote   string    `json:"c;osed_note"`
	UserModified int       `json:"user_modified"`
	TplCreator   string    `json:"tpl_creator"`
	ExpressionId int       `json:"expression_id"`
	StrategyId   int       `json:"strategy_id"`
	TemplateId   int       `json:"template_id"`
	Evnets       []*Events `json:"evevnts" orm:"reverse(many)"`
}

type Events struct {
	Id          int         `json:"id" orm:"pk"`
	Step        int         `json:"step"`
	Cond        string      `json:"cond"`
	Timestamp   time.Time   `json:"timestamp"`
	EventCaseId *EventCases `json:"event_caseId" orm:"rel(fk)"`
}

func insertEvent(q orm.Ormer, eve *coommonModel.Event) (res interface{}, err error) {
	sqltemplete := `INSERT INTO events (
		event_caseId,
		step,
		cond,
		timestamp
	) VALUES(?,?,?,?)`
	res, err = q.Raw(
		sqltemplete,
		eve.Id,
		eve.CurrentStep,
		fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
		time.Unix(eve.EventTime, 0),
	).Exec()
	return
}
func InsertEvent(eve *coommonModel.Event) {
	log := logger.Logger()
	q := orm.NewOrm()
	q.Using("falcon_portal")
	var event []EventCases
	q.Raw("select * from event_cases where id = ?", eve.Id).QueryRows(&event)
	if len(event) == 0 {
		//create cases
		sqltemplete := `INSERT INTO event_cases (
					id,
					endpoint,
					metric,
					func,
					cond,
					note,
					max_step,
					current_step,
					priority,
					status,
					timestamp,
					update_at,
					tpl_creator,
					expression_id,
					strategy_id,
					template_id
					) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
		res1, err := q.Raw(
			sqltemplete,
			eve.Id,
			eve.Endpoint,
			counterGen(eve.Metric(), utils.SortedTags(eve.PushedTags)),
			eve.Func(),
			//cond
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Strategy.Note,
			eve.MaxStep(),
			eve.CurrentStep,
			eve.Priority(),
			eve.Status,
			//start_at
			time.Unix(eve.EventTime, 0),
			//update_at
			time.Unix(eve.EventTime, 0),
			eve.Strategy.Tpl.Creator,
			eve.ExpressionId(),
			eve.StrategyId(),
			//template_id
			eve.TplId()).Exec()
		log.Debug(fmt.Printf("%v, %v", res1, err))

		//insert case
		res2, err := insertEvent(q, eve)
		log.Debug(fmt.Printf("%v, %v", res2, err))

	} else {
		//update cases
		res1, err := q.Raw(
			"UPDATE event_cases SET update_at = ?, current_step = ?, cond = ?, status = ? WHERE id = ?",
			time.Unix(eve.EventTime, 0),
			eve.CurrentStep,
			fmt.Sprintf("%v %v %v", eve.LeftValue, eve.Operator(), eve.RightValue()),
			eve.Status,
			eve.Id).Exec()
		log.Debug(fmt.Printf("%v, %v", res1, err))
		//insert case
		res2, err := insertEvent(q, eve)
		log.Debug(fmt.Printf("%v, %v", res2, err))
		log.Info("Hello")
	}
}

func counterGen(metric string, tags string) (mycounter string) {
	mycounter = metric
	if tags != "" {
		mycounter = fmt.Sprintf("%s/%s", metric, tags)
	}
	return
}
