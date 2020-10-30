package bean

import "github.com/go-xorm/xorm"

var (
	DbEngine *xorm.Engine
)

type DbHandler struct{}