package gocrud

import (
	"errors"
	"github.com/kordar/goi18n"
)

type ResourceService interface {
	ResourceName() string
	Search(body SearchBody) SearchVO
	SearchOne(body SearchBody) SearchOneVO
	Delete(body RemoveBody) error
	Add(body FormBody) (interface{}, error)
	Update(body FormBody) (interface{}, error)
	Edit(body EditorBody) error
	Configs(lang string) map[string]interface{}
}

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

func (mgr *ResourceManager) SelectOne(apiName string, searchBody SearchBody) (SearchOneVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(searchBody.Lang(), "errors", "resource.does.not.exist")
		return SearchOneVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).SearchOne(searchBody), nil
}

func (mgr *ResourceManager) Select(apiName string, searchBody SearchBody) (SearchVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(searchBody.Lang(), "errors", "resource.does.not.exist")
		return SearchVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Search(searchBody), nil
}

func (mgr *ResourceManager) Add(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(formBody.Lang(), "errors", "resource.does.not.exist")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Add(formBody)
}

func (mgr *ResourceManager) Update(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(formBody.Lang(), "errors", "resource.does.not.exist")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Update(formBody)
}

func (mgr *ResourceManager) Delete(apiName string, removeBody RemoveBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(removeBody.Lang(), "errors", "resource.does.not.exist")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Delete(removeBody)
}

func (mgr *ResourceManager) Edit(apiName string, editorBody EditorBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(editorBody.Lang(), "errors", "resource.does.not.exist")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Edit(editorBody)
}

func (mgr *ResourceManager) Configs(apiName string, lang string) (map[string]interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := goi18n.GetSectionValue(lang, "errors", "resource.does.not.exist")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Configs(lang), nil
}
