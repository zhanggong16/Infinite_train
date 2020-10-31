package request

import "github.com/labstack/echo"

//CommonContext is common context for all restful request
type CommonContext struct {
	RequestID  string
	Token      string
	TenantID   string
	TenantName string
	Pin        string
	Region     string
	Password   string
	IsAdmin    bool
}

type CustomContext struct {
	echo.Context
	CommonContext *CommonContext
}