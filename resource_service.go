package gocrud

import (
	"context"
	"errors"
)

type ResourceService interface {
	ResourceName() string
	DriverName() string
	Configs(ctx context.Context) map[string]interface{}
	Search(body SearchBody) SearchVO
	SearchOne(body SearchBody) SearchOneVO
	Remove(body RemoveBody) error
	Create(body FormBody) (interface{}, error)
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

func (common *CommonResourceService) Remove(body RemoveBody) error {
	message := MessageFn(body.Ctx(), "remove not found")
	return errors.New(message)
}

func (common *CommonResourceService) Create(body FormBody) (interface{}, error) {
	message := MessageFn(body.Ctx(), "create not found")
	return nil, errors.New(message)
}

func (common *CommonResourceService) Update(body FormBody) (interface{}, error) {
	message := MessageFn(body.Ctx(), "update not found")
	return nil, errors.New(message)
}

func (common *CommonResourceService) Edit(body EditorBody) error {
	message := MessageFn(body.Ctx(), "edit not found")
	return errors.New(message)
}

func (common *CommonResourceService) DriverName() string {
	return "gorm"
}

func (common *CommonResourceService) Configs(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{}
}
