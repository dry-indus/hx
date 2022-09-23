package common

type Page struct {
	PageNumber int64 `binding:"required,min=1"`        // 当前页
	PageSize   int64 `binding:"required,min=1,max=20"` // 每页条数
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