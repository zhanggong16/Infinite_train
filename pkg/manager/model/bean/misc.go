package bean

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"errors"
	"strconv"
)

func (bn *DbHandler) GetInstanceRecordByUuid(requestID string, bpId string) (int64, error) {
	result, err := DbEngine.Query("select sum(storage_size) as sum from backups where plan_id = ?", bpId)
	if err != nil {
		golog.Errorf(requestID, "select DB error %s", err.Error())
		return 0, err
	}
	if result == nil {
		err2 := errors.New("GetInstanceRecordByUuid result is empty")
		return 0, err2
	}
	resultInt, _ := strconv.ParseInt(string(result[0]["sum"]), 10, 64)
	return resultInt, nil
}