package parse

import (
	"github.com/lios/go-activiti/engine/impl/bpmn"
	. "github.com/lios/go-activiti/engine/impl/bpmn/model"
	. "github.com/lios/go-activiti/engine/impl/persistence/entity"
)

type ProcessParseHandler struct {
	AbstractActivityBpmnParseHandler
}

func (processParseHandler ProcessParseHandler) GetHandledType() string {
	return Process{}.GetType()
}

func (processParseHandler ProcessParseHandler) ExecuteParse(bpmnParse BpmnParse, flow bpmn.BaseElement) {
	bpmnParse.ProcessDefinitions = append(bpmnParse.ProcessDefinitions, transformProcess(bpmnParse, flow))
}

func transformProcess(bpmnParse BpmnParse, flow bpmn.BaseElement) ProcessDefinitionEntity {
	definitionEntityImpl := ProcessDefinitionEntityImpl{}
	process := flow.(*Process)
	definitionEntityImpl.Key = process.GetId()
	bpmnParse.ProcessFlowElements(process.FlowElementList)
	return &definitionEntityImpl
}
