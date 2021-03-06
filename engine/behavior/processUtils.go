package behavior

import (
	. "github.com/lios/go-activiti/engine"
	. "github.com/lios/go-activiti/engine/manager"
	"github.com/lios/go-activiti/errs"
	"sync"
)

var flowMap *sync.Map

type ProcessUtils struct {
	ProcessId int64
}

func init() {
	flowMap = new(sync.Map)
}
func (processUtil *ProcessUtils) GetCurrentTask(taskId int) (FlowElement, error) {
	manager := GetTaskManager()
	task, err := manager.FindById(taskId)
	if err != nil {
		return nil, err
	}
	processUtil.ProcessId = task.ProcessInstanceId
	processUtil.LoadProcess()
	defineManager := GetDefineManager()
	bytearry, err := defineManager.FindProcessByTask(task.ProcessInstanceId)
	if err != nil {
		return nil, err
	}
	currentTask := FindCurrentTask(bytearry, task.TaskDefineKey)
	return currentTask, nil
}

func (processUtil *ProcessUtils) GetFlowElement(flowElementId string) (FlowElement, error) {
	flowM, ok := flowMap.Load(processUtil.ProcessId)
	if !ok {

	}
	flowElements, isOk := flowM.(map[string]FlowElement)
	if isOk {
		flowElement := flowElements[flowElementId]
		return flowElement, nil
	}
	return nil, errs.ProcessError{Code: "1003", Msg: "not find"}
}

func (processUtil *ProcessUtils) LoadProcess() error {
	defineManager := GetDefineManager()
	bytearry, err := defineManager.FindProcessByTask(processUtil.ProcessId)
	if err != nil {
		return err
	}
	process := ConverterBpmn(bytearry)
	flowMap.Store(processUtil.ProcessId, process.FlowMap)
	return nil
}
