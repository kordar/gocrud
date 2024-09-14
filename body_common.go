package gocrud

import "context"

type CommonBody struct {
	ctx    context.Context
	driver string
}

func NewCommonBody(driverName string, ctx context.Context) *CommonBody {
	return &CommonBody{ctx: ctx, driver: driverName}
}

func (common *CommonBody) DriverName(params map[string]string) string {
	if params == nil || params["driver"] == "" {
		return common.driver
	}
	return params["driver"]
}

func (common *CommonBody) Ctx() context.Context {
	return common.ctx
}

func (common *CommonBody) LoadDriverName(params map[string]string) map[string]string {
	if params == nil {
		params = map[string]string{}
	}
	params["driver"] = common.DriverName(params)
	return params
}
