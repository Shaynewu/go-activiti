package entity

import (
	. "github.com/lios/go-activiti/engine/impl/persistence/entity/data"
)

var resourceDataManager ResourceDataManager

type ResourceEntityManagerImpl struct {
	AbstractEntityManager
}

func (resourceEntityManager ResourceEntityManagerImpl) GetDataManager() DataManagers {
	return resourceDataManager
}

func (resourceEntityManager ResourceEntityManagerImpl) FindResourcesByDeploymentId(deploymentId int64) ResourceEntity {
	resource, err := resourceDataManager.FindDeployedProcessDefinitionByDeploymentId(deploymentId)
	if err != nil {
		panic(err)
	}
	resourceEntityImpl := ResourceEntityImpl{}
	resourceEntityImpl.SetName(resource.Name)
	resourceEntityImpl.SetDeploymentId(resource.DeploymentId)
	resourceEntityImpl.SetBytes([]byte(resource.Bytes))
	return &resourceEntityImpl
}
