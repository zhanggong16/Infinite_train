package utils

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"fmt"
	"github.com/go-xorm/xorm"
)

func CreateOrmEngine(account string, password string, ip string, port int, schema string,
	charset string, maxIdle int, maxOpen int) (engine *xorm.Engine, err error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", account, password,
		ip, port, schema, charset)
	golog.Infof("0", "CreateOrmEngine using data source:%s", dataSourceName)
	engine, err = xorm.NewEngine("mysql", dataSourceName)
	if err == nil {
		engine.SetMaxIdleConns(maxIdle)
		engine.SetMaxOpenConns(maxOpen)

	}
	golog.Infof("0", "CreateOrmEngine finish")
	return engine, err
}