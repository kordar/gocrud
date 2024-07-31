package gocrud

type CommonBody struct {
	ctx    interface{}
	driver string
}

func NewCommonBody(driverName string, ctx interface{}) *CommonBody {
	return &CommonBody{ctx: ctx, driver: driverName}
}

func (common *CommonBody) DriverName(parallel map[string]string) string {
	if parallel == nil || parallel["driver"] == "" {
		return common.driver
	}
	return parallel["driver"]
}

func (common *CommonBody) Ctx() interface{} {
	return common.ctx
}

func (common *CommonBody) LoadDriverName(parallel map[string]string) map[string]string {
	if parallel == nil {
		parallel = map[string]string{}
	}
	parallel["driver"] = common.DriverName(parallel)
	return parallel
}
