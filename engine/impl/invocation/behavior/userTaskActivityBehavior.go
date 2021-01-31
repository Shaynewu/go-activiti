package behavior

import (
	. "github.com/lios/go-activiti/engine/common"
	"github.com/lios/go-activiti/engine/impl/bpmn/model"
	"github.com/lios/go-activiti/engine/impl/context"
	. "github.com/lios/go-activiti/engine/impl/handler"
	"github.com/lios/go-activiti/engine/impl/invocation"
	"github.com/lios/go-activiti/engine/impl/manager"
	"github.com/lios/go-activiti/engine/impl/persistence/entity"
	"github.com/lios/go-activiti/event"
	. "github.com/lios/go-activiti/event/impl"
	. "github.com/lios/go-activiti/model"
	"time"
)

type UserTaskActivityBehavior struct {
	UserTask   model.UserTask
	ProcessKey string
}

//普通用户节点处理
func (user UserTaskActivityBehavior) Execute(execution entity.ExecutionEntity) (err error) {
	task := Task{}
	task.ProcessInstanceId = execution.GetProcessInstanceId()
	task.Assignee = user.UserTask.Assignee
	task.StartTime = time.Now()
	task.TaskDefineKey = user.UserTask.Id
	task.TaskDefineName = user.UserTask.Name
	taskManager := manager.GetDataManager().TaskDataManager
	taskManager.Task = &task
	err = taskManager.Insert(execution)
	if err != nil {
		return err
	}
	err = handleAssignments(user.UserTask, task.Id, task.ProcessInstanceId)

	// All properties set, now firing 'create' events
	if interceptor.GetProcessEngineConfiguration().EventDispatcher.IsEnabled() {
		activitiEntityEvent, err := CreateEntityEvent(event.TASK_CREATED, task)
		if err != nil {
			return err
		}
		interceptor.GetProcessEngineConfiguration().EventDispatcher.DispatchEvent(activitiEntityEvent)
	}
	extensionElements := user.UserTask.ExtensionElements
	if extensionElements.TaskListener != nil && len(extensionElements.TaskListener) > 0 {
		taskListeners := extensionElements.TaskListener
		for _, listener := range taskListeners {
			if listener.EventType == TASK_TYPE_CREATE {
				err = PerformTaskListener(execution, user.UserTask.Name, user.ProcessKey)
				if err != nil {
					return err
				}
			}
		}
	}
	return err
}

//保存候选用户
func handleAssignments(user model.UserTask, taskId, processInstanceId int64) (err error) {
	users := user.CandidateUsers
	if len(users) >= 0 {
		for _, user := range users {
			link := IdentityLink{}
			link.TaskId = taskId
			link.ProcessInstanceId = processInstanceId
			link.UserId = user
			identityLinkManager := manager.GetDataManager().IdentityLinkDataManager
			identityLinkManager.IdentityLink = link
			err = identityLinkManager.CreateIdentityLink()
			if err != nil {
				return err
			}
		}
	}
	return err
}

//普通用户节点处理
func (user UserTaskActivityBehavior) Trigger(execution entity.ExecutionEntity) {
	user.Leave(execution)
}

func (user UserTaskActivityBehavior) Leave(execution entity.ExecutionEntity) {
	element := execution.GetCurrentFlowElement()
	execution.SetCurrentFlowElement(element)
	invocation.GetAgenda().PlanTakeOutgoingSequenceFlowsOperation(execution, true)
}