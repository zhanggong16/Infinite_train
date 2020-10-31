package restful

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/api/view"
	"github.com/labstack/echo"
	"net/http"
)

func (s *Server) Instances(c echo.Context) error {
	cc := c.(*request.CustomContext)
	requestID := cc.CommonContext.RequestID
	tenantId := cc.CommonContext.TenantID
	gid := c.Param("gid")

	golog.Infof(requestID, "Call Instances api by gid [%s], tenant id [%s], message %+v", gid, tenantId, cc.CommonContext)

	err := s.Validate.Var(gid, "required")
	if err != nil {
		errorResp := view.NewResponseError(constant.ParamErrorCode, requestID, err.Error())
		return c.JSON(errorResp.Code.GetHTTPCode(), errorResp.GetResponseErrorBody())
	}

	resultView := &view.CommonGidView{Gid: gid}
	ResponseView := &view.CommonSUCView{RequestId: requestID, Result: resultView}
	return c.JSON(http.StatusOK, ResponseView)
}