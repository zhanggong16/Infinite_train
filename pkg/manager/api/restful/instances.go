package restful

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
	"Infinite_train/pkg/manager/controller"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) InstancesDescribe(c echo.Context) error {
	ctx := c.(*request.CustomContext)
	requestID := ctx.CommonContext.RequestID
	tenantID := ctx.CommonContext.TenantID
	gid := c.Param("gid")

	golog.Infof(requestID, "instances describe by gid [%s], tenant id [%s], message %+v", gid, tenantID, ctx.CommonContext)
	// validate
	err := s.Validate.Var(gid, "required")
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	instanceID := gid
	// call controller
	resultView, errorResp := controller.InstancesControllerImpl.GetInstances(ctx.CommonContext, instanceID)
	if errorResp != nil {
		golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp)
	}
	// return view
	ResponseView := &view.CommonSUCView{RequestID: requestID, Result: resultView}
	return c.JSON(http.StatusOK, ResponseView)
}

func (s *Server) InstanceCreate(c echo.Context) error {
	ctx := c.(*request.CustomContext)
	requestID := ctx.CommonContext.RequestID
	//tenantID := ctx.CommonContext.TenantID

	requestBody := &request.InstanceCreateRequestBody{}
	err := c.Bind(requestBody)
	if err != nil {
		errorResp := view.NewResponseError(constant.BindParamErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "pin: %s, error: %s", ctx.CommonContext.Pin, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	err = s.Validate.Struct(requestBody)
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "pin: %s, error: %s", ctx.CommonContext.Pin, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	golog.Infof(requestID, "instance create for gid [%s]", requestBody.GID)
	// call controller

	// return view
	resultView := &view.InstanceDetailView{GID: requestBody.GID, Name: requestBody.Name}
	ResponseView := &view.CommonSUCView{RequestID: requestID, Result: resultView}
	return c.JSON(http.StatusOK, ResponseView)
}

func (s *Server) InstanceAction(c echo.Context) error {
	ctx := c.(*request.CustomContext)
	requestID := ctx.CommonContext.RequestID
	tenantID := ctx.CommonContext.TenantID
	pin := ctx.CommonContext.Pin
	gid := c.Param("gid")
	instancesController := controller.InstancesControllerImpl

	// validate
	err := s.Validate.Var(gid, "required")
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	requestBody := &request.InstanceActionBody{}
	err = c.Bind(requestBody)
	if err != nil {
		errorResp := view.NewResponseError(constant.BindParamErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "pin: %s, error: %s", pin, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	err = s.Validate.Struct(requestBody)
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		golog.Errorf(requestID, "pin: %s, error: %s", pin, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	golog.Infof(requestID, "instance action for gid [%s]", gid)
	instanceID := gid

	switch requestBody.Actions.Method {
	case "modify_name":
		newInstanceName := ""
		managerCommonContext := &request.ManagerCommonContext{ID: instanceID, CommonContext: ctx.CommonContext}
		if requestBody.Actions.Params != nil {
			bodyParam := &request.ChangeClusterNameParam{}
			var param []byte
			param, err = json.Marshal(requestBody.Actions.Params)
			if err != nil {
				errorResp := view.NewResponseError(constant.BindParamErrorCode, requestID, err.Error())
				golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
				return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
			}
			err = json.Unmarshal(param, bodyParam)
			if err != nil {
				errorResp := view.NewResponseError(constant.BindParamErrorCode, requestID, err.Error())
				golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
				return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
			}
			err = s.Validate.Struct(bodyParam)
			if err != nil {
				errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
				golog.Errorf(requestID, "tenant_id: %s, error message: %s", pin, errorResp.Error())
				return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
			}
			newInstanceName = bodyParam.NewName
		}
		// call controller
		errorResp := instancesController.ChangeInstanceName(managerCommonContext, newInstanceName)
		if errorResp != nil {
			return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
		}

		resultView := &view.InstanceDetailView{GID: gid, Name: newInstanceName}
		ResponseView := &view.CommonSUCView{RequestID: requestID, Result: resultView}
		return c.JSON(http.StatusAccepted, ResponseView)
	default:
		errorResp := view.NewResponseError(constant.StatusMethodNotAllowedErrorCode, requestBody.Actions.Method)
		return c.JSON(http.StatusBadRequest, errorResp.GetResponseErrorBody())
	}
}