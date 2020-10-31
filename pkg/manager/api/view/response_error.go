package view

import (
	"Infinite_train/pkg/common/constant"
	"fmt"
)

//ResponseError resp
type ResponseError struct {
	Code      constant.StatusCode `json:"code"`
	Message   string              `json:"message"`
	Status    string              `json:"status"`
	RequestId string              `json:"requestId"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("code:[%v]:description[%v]", e.Code, e.Message)
}

//NewResponseError create ResponseError
func NewResponseError(code constant.StatusCode, desciption ...interface{}) *ResponseError {
	msg := constant.ErrorMessage(code)
	if desciption != nil {
		return &ResponseError{Code: code, Message: fmt.Sprintf(msg, desciption...), Status: "ERROR", RequestId: desciption[0].(string)}
	}
	return &ResponseError{Code: code, Message: msg, Status: "ERROR", RequestId: ""}
}

type ResponseErrorBody struct {
	RequestId string         `json:"requestId"`
	Error     *ResponseError `json:"error"`
}

func (e *ResponseError) GetResponseErrorBody() *ResponseErrorBody {
	return &ResponseErrorBody{RequestId: e.RequestId, Error: e}
}