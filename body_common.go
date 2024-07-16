package gocrud

type CommonBody[C interface{}] struct {
	ctx    C
	driver string
}

func NewCommonBody[C interface{}](driverName string, ctx C) *CommonBody[C] {
	return &CommonBody[C]{ctx: ctx, driver: driverName}
}

func (common *CommonBody[C]) DriverName(parallel map[string]string) string {
	if parallel == nil || parallel["driver"] == "" {
		return common.driver
	}
	return parallel["driver"]
}

func (common *CommonBody[C]) Ctx() C {
	return common.ctx
}

func (common *CommonBody[C]) LoadDriverName(parallel map[string]string) map[string]string {
	if parallel == nil {
		parallel = map[string]string{}
	}
	parallel["driver"] = common.DriverName(parallel)
	return parallel
}
