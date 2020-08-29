package impl

import . "github.com/lios/go-activiti/event"

type ActivitiEntityEventImpl struct {
	ActivitiEntityEvent
	ActivitiEventImpl
	Entity interface{}
}

func (ActivitiEntityEventImpl) GetType() ActivitiEventType {
	return TASK_CREATED
}
