package common

type Page struct {
	//Limit (min 1)
	PageNumber int64 `json:"pageNumber" binding:"required,min=1" validate:"required"` // 当前页
	//Limit (min 1,max 20)
	PageSize int64 `json:"pageSize" binding:"required,min=1,max=20" validate:"required"` // 每页条数
}

func (p *Page) Skip() int64 {
	return (p.PageNumber - 1) * p.PageSize
}

func (p *Page) Limit() int64 {
	return p.PageSize
}

type Logger interface {
	Errorf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}
