package gocrud

type SearchVO struct {
	Data  interface{} `json:"data,omitempty"`
	Count int64       `json:"count"`
}

func NewSearchVO(data interface{}, count int64) SearchVO {
	return SearchVO{
		Data:  data,
		Count: count,
	}
}

type SearchOneVO struct {
	Info interface{} `json:"info,omitempty"`
}

func NewSearchOneVO(info interface{}) SearchOneVO {
	return SearchOneVO{
		Info: info,
	}
}
