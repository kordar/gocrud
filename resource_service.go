package gocrud

import (
	"errors"
	"github.com/kordar/gocfg"
)

type ResourceService interface {
	ResourceName() string
	DriverName() string
	Configs(ctx interface{}) map[string]interface{}
	Search(body SearchBody) SearchVO
	SearchOne(body SearchBody) SearchOneVO
	Delete(body RemoveBody) error
	Add(body FormBody) (interface{}, error)
	Update(body FormBody) (interface{}, error)
	Edit(body EditorBody) error
}

type CommonResourceService struct {
}

func (common *CommonResourceService) Search(body SearchBody) SearchVO {
	return SearchVO{}
}

func (common *CommonResourceService) SearchOne(body SearchBody) SearchOneVO {
	return SearchOneVO{}
}

func (common *CommonResourceService) Delete(body RemoveBody) error {
	message := gocfg.GetSectionValue(LangFn(body.Ctx()), "resource.errors.no_provided", "language")
	return errors.New(message)
}

func (common *CommonResourceService) Add(body FormBody) (interface{}, error) {
	message := gocfg.GetSectionValue(LangFn(body.Ctx()), "resource.errors.no_provided", "language")
	return nil, errors.New(message)
}

func (common *CommonResourceService) Update(body FormBody) (interface{}, error) {
	message := gocfg.GetSectionValue(LangFn(body.Ctx()), "resource.errors.no_provided", "language")
	return nil, errors.New(message)
}

func (common *CommonResourceService) Edit(body EditorBody) error {
	message := gocfg.GetSectionValue(LangFn(body.Ctx()), "resource.errors.no_provided", "language")
	return errors.New(message)
}

func (common *CommonResourceService) DriverName() string {
	return "gorm"
}

func (common *CommonResourceService) Configs(ctx interface{}) map[string]interface{} {
	return map[string]interface{}{}
}
