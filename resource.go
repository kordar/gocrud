package gocrud

import (
	"errors"
	"github.com/kordar/gocfg"
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
	Driver() string
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
		message := gocfg.GetSectionValue(lang(), "resource.errors.view_not_exist", "language")
		return SearchOneVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).SearchOne(searchBody), nil
}

func (mgr *ResourceManager) Select(apiName string, searchBody SearchBody) (SearchVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.list_not_exist", "language")
		return SearchVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Search(searchBody), nil
}

func (mgr *ResourceManager) Add(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.add_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Add(formBody)
}

func (mgr *ResourceManager) Update(apiName string, formBody FormBody) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.update_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Update(formBody)
}

func (mgr *ResourceManager) Delete(apiName string, removeBody RemoveBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.delete_not_exist", "language")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Delete(removeBody)
}

func (mgr *ResourceManager) Edit(apiName string, editorBody EditorBody) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.edit_not_exist", "language")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Edit(editorBody)
}

func (mgr *ResourceManager) Configs(apiName string) (map[string]interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.configs_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Configs(lang()), nil
}

func (mgr *ResourceManager) Driver(apiName string) string {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		message := gocfg.GetSectionValue(lang(), "resource.errors.driver_name_required", "language")
		return message
	}
	return mgr.container.GetResourceService(apiName).Driver()
}
