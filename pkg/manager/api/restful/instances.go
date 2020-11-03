package restful

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
	"Infinite_train/pkg/manager/controller"
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) DescribeInstances(c echo.Context) error {
	ctx := c.(*request.CustomContext)
	requestID := ctx.CommonContext.RequestID
	tenantID := ctx.CommonContext.TenantID
	gid := c.Param("gid")

	golog.Infof(requestID, "Call Instances api by gid [%s], tenant id [%s], message %+v", gid, tenantID, ctx.CommonContext)
	// validate
	err := s.Validate.Var(gid, "required")
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}
	// call controller
	resultView, errorResp := controller.InstancesControllerImpl.GetInstances(ctx, gid)
	if errorResp != nil {
		golog.Errorf(requestID, "tenant_id: %s, error message: %s", tenantID, errorResp.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp)
	}
	// return view
	ResponseView := &view.CommonSUCView{RequestId: requestID, Result: resultView}
	return c.JSON(http.StatusOK, ResponseView)
}