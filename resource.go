package gocrud

import (
	"context"
	"errors"
)

type resourceContainer struct {
	resourceServiceHashMap map[string]ResourceService
}

func (container *resourceContainer) AddResourceService(service ResourceService) {
	container.resourceServiceHashMap[service.ResourceName()] = service
}

func (container *resourceContainer) GetResourceService(apiName string) ResourceService {
	return container.resourceServiceHashMap[apiName]
}

type ResourceManager struct {
	container resourceContainer
}

func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		container: resourceContainer{
			resourceServiceHashMap: make(map[string]ResourceService),
		},
	}
}

func (mgr *ResourceManager) AddResourceService(service ResourceService) {
	mgr.container.AddResourceService(service)
}

func (mgr *ResourceManager) ReadOne(apiName string, searchBody SearchBody) (SearchOneVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(searchBody.Ctx(), "read-one not found")
		return SearchOneVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).SearchOne(searchBody), nil
}

func (mgr *ResourceManager) Read(apiName string, searchBody SearchBody) (SearchVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(searchBody.Ctx(), "read not found")
		return SearchVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Search(searchBody), nil
}

func (mgr *ResourceManager) Create(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(formBody.Ctx(), "create not found")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Create(formBody)
}

func (mgr *ResourceManager) Update(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(formBody.Ctx(), "update not found")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Update(formBody)
}

func (mgr *ResourceManager) Delete(apiName string, removeBody RemoveBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(removeBody.Ctx(), "remove not found")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Remove(removeBody)
}

func (mgr *ResourceManager) Edit(apiName string, editorBody EditorBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(editorBody.Ctx(), "edit not found")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Edit(editorBody)
}

func (mgr *ResourceManager) Configs(apiName string, ctx context.Context) (map[string]interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(ctx, "create not found")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Configs(ctx), nil
}

func (mgr *ResourceManager) DriverName(apiName string, ctx context.Context) string {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := MessageFn(ctx, "driver name required")
		return message
	}
	return mgr.container.GetResourceService(apiName).DriverName()
}
