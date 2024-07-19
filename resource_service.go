package gocrud

import (
	"errors"
	"github.com/kordar/gocfg"
)

type ResourceService[T interface{}, C interface{}] interface {
	ResourceName() string
	DriverName() string
	Configs() map[string]interface{}
	Search(body SearchBody[T, C]) SearchVO
	SearchOne(body SearchBody[T, C]) SearchOneVO
	Delete(body RemoveBody[T, C]) error
	Add(body FormBody[T, C]) (interface{}, error)
	Update(body FormBody[T, C]) (interface{}, error)
	Edit(body EditorBody[T, C]) error
}

type CommonResourceService[T interface{}, C interface{}] struct {
}

func (common *CommonResourceService[T, C]) Search(body SearchBody[T, C]) SearchVO {
	return SearchVO{}
}

func (common *CommonResourceService[T, C]) SearchOne(body SearchBody[T, C]) SearchOneVO {
	return SearchOneVO{}
}

func (common *CommonResourceService[T, C]) Delete(body RemoveBody[T, C]) error {
	message := gocfg.GetSectionValue(Lang(body.Ctx()), "resource.errors.no_provided", "language")
	return errors.New(message)
}

func (common *CommonResourceService[T, C]) Add(body FormBody[T, C]) (interface{}, error) {
	message := gocfg.GetSectionValue(Lang(body.Ctx()), "resource.errors.no_provided", "language")
	return nil, errors.New(message)
}

func (common *CommonResourceService[T, C]) Update(body FormBody[T, C]) (interface{}, error) {
	message := gocfg.GetSectionValue(Lang(body.Ctx()), "resource.errors.no_provided", "language")
	return nil, errors.New(message)
}

func (common *CommonResourceService[T, C]) Edit(body EditorBody[T, C]) error {
	message := gocfg.GetSectionValue(Lang(body.Ctx()), "resource.errors.no_provided", "language")
	return errors.New(message)
}

func (common *CommonResourceService[T, C]) DriverName() string {
	return "gorm"
}

func (common *CommonResourceService[T, C]) Configs() map[string]interface{} {
	return map[string]interface{}{}
}
