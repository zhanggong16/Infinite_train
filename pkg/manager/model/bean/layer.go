package bean

import "github.com/go-xorm/xorm"

type DBBeanInterface interface {
	InsertRecord(requestId, tableName string, keyValues map[string]interface{}) error
	UpdateRecord(requestId, tableName string, setKeyValue, condition map[string]interface{}) (int64, error)
	UpdateRecordIn(requestId, tableName string, setKeyValue map[string]interface{}, inKey string, inValue []string) (int64, error)
	GetRecord(requestId, tableName string, condition map[string]interface{}) (interface{}, error)
	UpdateRecordWhere(requestId, tableName string, setKeyValue map[string]interface{}, query interface{}, args []interface{}) (int64, error)
	FindRecords(requestId, tableName string, condition map[string]interface{}, rowsSlices interface{}) error
	FindRecordsWhere(requestId, tableName string, query interface{}, args []interface{}, rowsSlices interface{}) error
	CountRecords(requestId, tableName string, condition map[string]interface{}) (int64, error)
	FindRecordsIn(requestId, tableName string, inKey string, inValue []string, rowsSlices interface{}) error
	FindRecordsWhereIn(requestId, tableName string, inKey string, inValue []string, rowsSlices interface{}) error
	FindRecordsDesc(requestId, tableName string, where map[string]interface{}, descName string, rowsSlices interface{}) error
	FindRecordsAscIn(requestId, tableName string, where map[string]interface{}, inKey string, inValue []string, ascName string, rowsSlices interface{}) error
	TxBegin() (*xorm.Session, error)
	TxRollback(session *xorm.Session) error
	TxCommit(session *xorm.Session) error
	TxInsertRecord(txSession *xorm.Session, requestId, tableName string, keyValues map[string]interface{}) error
	TxUpdateRecord(txSession *xorm.Session, requestId, tableName string, setKeyValue, condition map[string]interface{}) (int64, error)
	FindRecordsForFilterGroups(requestId, tableName string, condition *FilterGroups, rowsSlices interface{}) (total int64, err error)
	DeleteRecord(requestId, tableName string, condition map[string]interface{}) (int64, error)
}

var DBBeanImpl DBBeanInterface = new(DbHandler)