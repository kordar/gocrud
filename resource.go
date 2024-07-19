package gocrud

import (
	"errors"
	"github.com/kordar/gocfg"
)

type resourceContainer[T interface{}, C interface{}] struct {
	resourceServiceHashMap map[string]ResourceService[T, C]
}

func (container *resourceContainer[T, C]) AddResourceService(service ResourceService[T, C]) {
	container.resourceServiceHashMap[service.ResourceName()] = service
}

func (container *resourceContainer[T, C]) GetResourceService(apiName string) ResourceService[T, C] {
	return container.resourceServiceHashMap[apiName]
}

type ResourceManager[T interface{}, C interface{}] struct {
	container resourceContainer[T, C]
}

func NewResourceManager[T interface{}, C interface{}]() *ResourceManager[T, C] {
	return &ResourceManager[T, C]{
		container: resourceContainer[T, C]{
			resourceServiceHashMap: make(map[string]ResourceService[T, C]),
		},
	}
}

func (mgr *ResourceManager[T, C]) AddResourceService(service ResourceService[T, C]) {
	mgr.container.AddResourceService(service)
}

func (mgr *ResourceManager[T, C]) SelectOne(apiName string, searchBody SearchBody[T, C]) (SearchOneVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(searchBody.Ctx())
		if lang == "" {
			return SearchOneVO{}, errors.New("resource view not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.view_not_exist", "language")
		return SearchOneVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).SearchOne(searchBody), nil
}

func (mgr *ResourceManager[T, C]) Select(apiName string, searchBody SearchBody[T, C]) (SearchVO, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(searchBody.Ctx())
		if lang == "" {
			return SearchVO{}, errors.New("resource list not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.list_not_exist", "language")
		return SearchVO{}, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Search(searchBody), nil
}

func (mgr *ResourceManager[T, C]) Add(apiName string, formBody FormBody[T, C]) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(formBody.Ctx())
		if lang == "" {
			return nil, errors.New("resource add not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.add_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Add(formBody)
}

func (mgr *ResourceManager[T, C]) Update(apiName string, formBody FormBody[T, C]) (interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(formBody.Ctx())
		if lang == "" {
			return nil, errors.New("resource update not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.update_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Update(formBody)
}

func (mgr *ResourceManager[T, C]) Delete(apiName string, removeBody RemoveBody[T, C]) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(removeBody.Ctx())
		if lang == "" {
			return errors.New("resource delete not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.delete_not_exist", "language")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Delete(removeBody)
}

func (mgr *ResourceManager[T, C]) Edit(apiName string, editorBody EditorBody[T, C]) error {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(editorBody.Ctx())
		if lang == "" {
			return errors.New("resource edit not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.edit_not_exist", "language")
		return errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Edit(editorBody)
}

func (mgr *ResourceManager[T, C]) Configs(apiName string, ctx interface{}) (map[string]interface{}, error) {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(ctx)
		if lang == "" {
			return nil, errors.New("resource configs not exist")
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.configs_not_exist", "language")
		return nil, errors.New(message)
	}
	return mgr.container.GetResourceService(apiName).Configs(), nil
}

func (mgr *ResourceManager[T, C]) DriverName(apiName string, ctx interface{}) string {
	if apiName == "" || mgr.container.GetResourceService(apiName) == nil {
		lang := Lang(ctx)
		if lang == "" {
			return "resource driver name required"
		}
		message := gocfg.GetSectionValue(lang, "resource.errors.driver_name_required", "language")
		return message
	}
	return mgr.container.GetResourceService(apiName).DriverName()
}
