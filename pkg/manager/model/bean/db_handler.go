package bean

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/model"
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"reflect"
	"strings"
)

var (
	DbEngine *xorm.Engine
)

type DbHandler struct{
}

//FiltersPara ...
type FiltersPara struct {
	Column   string
	Operator string
	Values   []interface{}
}

type SortsPara struct {
	Column    string
	Direction string
}

//Filter ...
type Filter struct {
	IsPaging  bool
	PageIndex int
	PageSize  int
	Filters   []*FiltersPara
	DescSorts []string
	AscSorts  []string
}

type FilterGroup struct {
	Filters []*FiltersPara
}

type FilterGroups struct {
	FilterGroup []*FilterGroup
	IsPaging    bool
	PageIndex   int
	PageSize    int
	BaseFilters []*FiltersPara
	DescSorts   []string
	AscSorts    []string
}

func snakeCasedName(name string) string {
	newStr := make([]rune, 0)
	for idx, chr := range name {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if idx > 0 {
				newStr = append(newStr, '_')
			}
			chr -= 'A' - 'a'
		}
		newStr = append(newStr, chr)
	}
	return string(newStr)
}

func convertQueryInfo(filters []*FiltersPara) (string, []interface{}) {
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for _, v := range filters {
		if len(v.Values) <= 0 {
			v.Values = []interface{}{""}
		}
		queryString = queryString + " and "
		// k reflect?
		//lt, le, gt, ge, ne, like,eq
		if v.Operator == "eq" {
			queryString = queryString + snakeCasedName(v.Column) + "=?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "lt" {
			queryString = queryString + snakeCasedName(v.Column) + "<?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "le" {
			queryString = queryString + snakeCasedName(v.Column) + "<=?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "gt" {
			queryString = queryString + snakeCasedName(v.Column) + ">?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "ge" {
			queryString = queryString + snakeCasedName(v.Column) + ">=?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "ne" {
			queryString = queryString + snakeCasedName(v.Column) + "<>?"
			queryArgs = append(queryArgs, v.Values[0])
		} else if v.Operator == "like" {
			queryString = queryString + "("
			likeArgs := []interface{}{}
			for _, item := range v.Values {
				val := fmt.Sprintf("%v", item)
				queryString = queryString + snakeCasedName(v.Column) + " like ?" + " or "
				if item == "" {
					likeArgs = append(likeArgs, val)
				} else {
					likeArgs = append(likeArgs, "%"+val+"%")
				}
			}
			queryString = strings.TrimSuffix(queryString, " or ")
			queryString = queryString + ")"
			queryArgs = append(queryArgs, likeArgs...)
		} else if v.Operator == "in" {
			queryString = queryString + snakeCasedName(v.Column) + " in ("
			inArgs := []interface{}{}
			for _, item := range v.Values {
				queryString = queryString + "?,"
				inArgs = append(inArgs, item)
			}
			queryString = strings.TrimSuffix(queryString, ",")
			queryString = queryString + ")"
			queryArgs = append(queryArgs, inArgs...)
		} else if v.Operator == "ni" {
			queryString = queryString + snakeCasedName(v.Column) + " not in ("
			inArgs := []interface{}{}
			for _, item := range v.Values {
				queryString = queryString + "?,"
				inArgs = append(inArgs, item)
			}
			queryString = strings.TrimSuffix(queryString, ",")
			queryString = queryString + ")"
			queryArgs = append(queryArgs, inArgs...)
		} else {
			queryString = queryString + snakeCasedName(v.Column) + "=?"
			queryArgs = append(queryArgs, v.Values[0])
		}
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	return queryString, queryArgs
}

func (bn *DbHandler) nameToStruct(tableStructName string) (interface{}, error) {
	switch tableStructName {
	case "Table1":
		return new(model.Table1), nil
	default:
		err := errors.New("invalid table name: " + tableStructName)
		return nil, err
	}
}

//InsertRecord instert a record to table
func (bn *DbHandler) InsertRecord(requestId, tableName string, keyValues map[string]interface{}) error {
	golog.Infof(requestId, "InsertRecord table:%s, key_values : %+v", tableName, keyValues)
	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	for k, v := range keyValues {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
	_, err := DbEngine.Insert(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) UpdateRecord(requestId, tableName string, setKeyValue, condition map[string]interface{}) (int64, error) {
	golog.Infof(requestId, "UpdateRecord table:%s setKeyValue:%+v condition: %+v ", tableName, setKeyValue, condition)
	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	for k, v := range setKeyValue {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		// k reflect?
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	result, err := DbEngine.Where(queryString, queryArgs...).Update(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return -1, err
	}
	return result, nil
}

func (bn *DbHandler) UpdateRecordIn(requestId, tableName string, setKeyValue map[string]interface{}, inKey string, inValue []string) (int64, error) {
	golog.Infof(requestId, "UpdateRecordIn table:%s setKeyValue:%+v inKey: %+v inValue: %+v", tableName, setKeyValue, inKey, inValue)
	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	for k, v := range setKeyValue {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
	result, err := DbEngine.In(snakeCasedName(inKey), inValue).Update(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return -1, err
	}
	if result <= 0 {
		errNew := errors.New("not effect any record")
		golog.Errorf(requestId, "error message : %s", errNew.Error())
		return result, errNew
	}
	return result, nil
}

func (bn *DbHandler) GetRecord(requestId, tableName string, condition map[string]interface{}) (interface{}, error) {
	golog.Infof(requestId, "GetRecord table:%s, condition:%+v ", tableName, condition)
	dbRecord, _ := bn.nameToStruct(tableName)
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	has, err := DbEngine.Where(queryString, queryArgs...).Get(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return nil, err
	}
	if !has {
		err2 := errors.New("query cluster get nothing")
		golog.Warnf(requestId, "In GetRecordWhere query error：%s", err2.Error())
		return nil, nil
	}
	return dbRecord, nil
}

func (bn *DbHandler) UpdateRecordWhere(requestId, tableName string, setKeyValue map[string]interface{}, query interface{}, args []interface{}) (int64, error) {
	golog.Infof(requestId, "UpdateRecordWhere table:%s setKeyValue:%+v query:%+v, args:%+v ", tableName,
		setKeyValue, query, args)
	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	for k, v := range setKeyValue {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
	result, err := DbEngine.Where(query, args...).Update(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return -1, err
	}
	if result <= 0 {
		errNew := errors.New("not effect any record")
		golog.Errorf(requestId, "error message : %s", errNew.Error())
		return result, errNew
	}
	return result, nil
}

func (bn *DbHandler) GetRecordWhere(requestId, tableName string, query interface{}, args []interface{}) (interface{}, error) {
	golog.Infof(requestId, "GetRecordWhere table:%s, query: %+v, args:%+v ", tableName, query, args)
	dbRecord, _ := bn.nameToStruct(tableName)
	has, err := DbEngine.Where(query, args...).Get(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return nil, err
	}
	if !has {
		err2 := errors.New("query cluster get nothing")
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err2.Error())
		return nil, err
	}
	return dbRecord, nil
}

func (bn *DbHandler) FindRecords(requestId, tableName string, condition map[string]interface{}, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecords table:%s, condition:%+v ", tableName, condition)
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	err := DbEngine.Where(queryString, queryArgs...).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) FindRecordsWhere(requestId, tableName string, query interface{}, args []interface{}, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecordsWhere table:%s, query: %+v, args:%+v ", tableName, query, args)
	err := DbEngine.Where(query, args...).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In FindRecordsWhere query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) CountRecords(requestId, tableName string, condition map[string]interface{}) (int64, error) {
	golog.Infof(requestId, "CountRecords table:%s, condition:%+v ", tableName, condition)
	dbRecord, _ := bn.nameToStruct(tableName)
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	total, err := DbEngine.Where(queryString, queryArgs...).Count(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return 0, err
	}
	return total, nil
}

func (bn *DbHandler) FindRecordsIn(requestId, tableName string, inKey string, inValue []string, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecordsIn table:%s, where Condition:%+v inKey:%+v inValue:+v",
		tableName, inKey, inValue)
	err := DbEngine.In(snakeCasedName(inKey), inValue).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) FindRecordsWhereIn(requestId, tableName string, inKey string, inValue []string, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecordsWhereIn table:%s, where Condition:%+v inKey:%+v inValue:+v",
		tableName, inKey, inValue)
	err := DbEngine.In(snakeCasedName(inKey), inValue).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In FindRecordsWhereIn query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) FindRecordsDesc(requestId, tableName string, where map[string]interface{}, descName string, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecordsDescIn table:%s, where Condition:%+v", tableName, where)
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range where {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	err := DbEngine.Desc(snakeCasedName(descName)).Where(queryString, queryArgs...).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) FindRecordsAscIn(requestId, tableName string, where map[string]interface{}, inKey string, inValue []string, ascName string, rowsSlices interface{}) error {
	golog.Infof(requestId, "FindRecordsAscIn table:%s, where Condition:%+v inKey:%+v inValue:+v", tableName, where, inKey, inValue)
	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range where {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	err := DbEngine.Asc(snakeCasedName(ascName)).Where(queryString, queryArgs...).In(snakeCasedName(inKey), inValue).Find(rowsSlices)
	if err != nil {
		golog.Errorf(requestId, "In GetRecordWhere query error：%s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) TxBegin() (*xorm.Session, error) {
	session := DbEngine.NewSession()
	err := session.Begin()
	return session, err
}

func (bn *DbHandler) TxRollback(session *xorm.Session) error {
	return session.Rollback()
}

func (bn *DbHandler) TxCommit(session *xorm.Session) error {
	return session.Commit()
}

func (bn *DbHandler) TxInsertRecord(txSession *xorm.Session, requestId, tableName string, keyValues map[string]interface{}) error {
	golog.Infof(requestId, "TxInsertRecord table:%s, key_values : %+v", tableName, keyValues)

	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	for k, v := range keyValues {
		val.FieldByName(k).Set(reflect.ValueOf(v))
	}
	_, err := txSession.Insert(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return err
	}
	return nil
}

func (bn *DbHandler) TxUpdateRecord(txSession *xorm.Session, requestId, tableName string, setKeyValue, condition map[string]interface{}) (int64, error) {
	golog.Infof(requestId, "TxUpdateRecord table:%s setKeyValue:%s condition: %+v ", tableName, setKeyValue, condition)

	dbRecord, _ := bn.nameToStruct(tableName)
	val := reflect.ValueOf(dbRecord).Elem()
	var needToUpdateColumns []string
	for k, v := range setKeyValue {
		val.FieldByName(k).Set(reflect.ValueOf(v))
		needToUpdateColumns = append(needToUpdateColumns, snakeCasedName(k))
	}

	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	result, err := txSession.Where(queryString, queryArgs...).Cols(needToUpdateColumns...).Update(dbRecord)
	if err != nil {
		golog.Errorf("", "error message : %s", err.Error())
		return -1, err
	}

	return result, nil
}

func (bn *DbHandler) FindRecordsForFilterGroups(requestId, tableName string, condition *FilterGroups, rowsSlices interface{}) (total int64, err error) {
	golog.Infof("", "FindRecordsForFilterGroups table:%s, condition:%+v ", tableName, condition)

	total = 0
	var err1 error

	//update cluster other fields
	dbRecord, _ := bn.nameToStruct(tableName)

	//update cluster other fields
	queryString := ""
	queryArgs := make([]interface{}, 0)
	baseString, baseArgs := convertQueryInfo(condition.BaseFilters)
	for _, filterGroup := range condition.FilterGroup {
		tempString, tempArgs := convertQueryInfo(filterGroup.Filters)
		if tempString == "" {
			continue
		}
		queryString = queryString + strings.Join([]string{tempString, "and "}, " ")
		queryArgs = append(queryArgs, tempArgs...)
	}
	queryString = strings.TrimSuffix(queryString, " and ")

	if queryString != "" {
		queryString = strings.Join([]string{baseString, "and", "(", queryString, ")"}, " ")
	} else {
		queryString = baseString
	}

	queryArgs = append(baseArgs, queryArgs...)

	golog.Infof(requestId, "query string %s", queryString)
	golog.Infof(requestId, "query args %v", queryArgs)

	total, err1 = DbEngine.Where(queryString, queryArgs...).Count(dbRecord)

	if condition.IsPaging {

		intTotal := int(total)
		queryPageSize := condition.PageSize
		queryPageIndex := condition.PageIndex

		var maxPageIndex int
		if intTotal == 0 {
			condition.PageIndex = 1
			return total, nil
		} else if intTotal%queryPageSize == 0 {
			maxPageIndex = intTotal / queryPageSize
		} else {
			maxPageIndex = intTotal/queryPageSize + 1
		}
		if maxPageIndex < condition.PageIndex {
			queryPageIndex = maxPageIndex
			condition.PageIndex = queryPageIndex
		}

		if len(condition.DescSorts) != 0 && len(condition.AscSorts) != 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Desc(condition.DescSorts...).
				Asc(condition.AscSorts...).
				Limit(condition.PageSize, (condition.PageIndex-1)*condition.PageSize).
				Find(rowsSlices)
		} else if len(condition.DescSorts) == 0 && len(condition.AscSorts) != 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Asc(condition.AscSorts...).
				Limit(condition.PageSize, (condition.PageIndex-1)*condition.PageSize).
				Find(rowsSlices)
		} else if len(condition.DescSorts) != 0 && len(condition.AscSorts) == 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Desc(condition.DescSorts...).
				Limit(condition.PageSize, (condition.PageIndex-1)*condition.PageSize).
				Find(rowsSlices)
		} else {
			err = DbEngine.Where(queryString, queryArgs...).
				Limit(condition.PageSize, (condition.PageIndex-1)*condition.PageSize).
				Find(rowsSlices)
		}

	} else {
		if len(condition.DescSorts) != 0 && len(condition.AscSorts) != 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Desc(condition.DescSorts...).
				Asc(condition.AscSorts...).
				Find(rowsSlices)
		} else if len(condition.DescSorts) == 0 && len(condition.AscSorts) != 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Asc(condition.AscSorts...).
				Find(rowsSlices)
		} else if len(condition.DescSorts) != 0 && len(condition.AscSorts) == 0 {
			err = DbEngine.Where(queryString, queryArgs...).
				Desc(condition.DescSorts...).
				Find(rowsSlices)
		} else {
			err = DbEngine.Where(queryString, queryArgs...).
				Find(rowsSlices)
		}
	}

	if err != nil {
		golog.Errorf(requestId, "In FindRecordsForFilterGroups query error：%s, error:%s", err.Error())
		return 0, err
	}

	if err1 != nil {
		golog.Errorf(requestId, "In FindRecordsForFilterGroups count error：%s, error:%s", err1.Error())
		return 0, err1
	}

	return total, nil
}

func (bn *DbHandler) DeleteRecord(requestId, tableName string, condition map[string]interface{}) (int64, error) {
	golog.Infof(requestId, "DeleteRecord table:%s condition: %+v ", tableName, condition)
	dbRecord, _ := bn.nameToStruct(tableName)

	queryString := ""
	queryArgs := make([]interface{}, 0)
	for k, v := range condition {
		// k reflect?
		queryString = queryString + " and " + snakeCasedName(k) + "=?"
		queryArgs = append(queryArgs, v)
	}
	queryString = strings.TrimPrefix(queryString, " and ")
	result, err := DbEngine.Where(queryString, queryArgs...).Delete(dbRecord)
	if err != nil {
		golog.Errorf(requestId, "error message : %s", err.Error())
		return -1, err
	}
	return result, nil
}