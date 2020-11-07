package constant

import "net/http"

type StatusCode int

/**
1000-1999 api error code
2000-3000 controller error code
*/
const (
	ParamErrorCode         StatusCode = 1000
	BindParamErrorCode     StatusCode = 1001
	JSONUnmarshalErrorCode StatusCode = 1002
	JSONMarshalErrorCode   StatusCode = 1003
	//2000 start
	InsertDBErrorCode                              StatusCode = 2000
	DeleteDBErrorCode                              StatusCode = 2001
	DeleteDBNotExistCode                           StatusCode = 2002
	SelectDBErrorCode                              StatusCode = 2003
	UpdateDBErrorCode                              StatusCode = 2004
	InstanceTaskStatusError                        StatusCode = 2005
	CreateFileError                                StatusCode = 2006
	PollingErrorCode                               StatusCode = 2007
	PollingSubTaskErrorCode                        StatusCode = 2008
	ScpFileErrorCode                               StatusCode = 2009
	PingAgentErrorCode                             StatusCode = 2010
	StatusMethodNotAllowedErrorCode			       StatusCode = 2011
)

var statusText = map[StatusCode]string{
	ParamErrorCode:          "RequestID %s parameter error [%s]",
	BindParamErrorCode:      "RequestID %s error message [%s]",
	InsertDBErrorCode:       "RequestID %s error message [%s]",
	DeleteDBErrorCode:       "RequestID %s error message [%s]",
	DeleteDBNotExistCode:    "RequestID %s error message [delete id: %s not exist]",
	SelectDBErrorCode:       "RequestID %s error message [%s]",
	UpdateDBErrorCode:       "RequestID %s error message [%s]",
	JSONUnmarshalErrorCode:  "RequestID %s error message [%s]",
	JSONMarshalErrorCode:    "RequestID %s error message [%s]",
	InstanceTaskStatusError: "RequestID %s select instance task ids [%s] not all success",
	CreateFileError:         "RequestID %s An error occurred with file opening or creation\n file [%s]",
	PollingErrorCode:        "RequestID %s ..polling...all ids %v status and success %v",
	PollingSubTaskErrorCode: "RequestID %s ..polling sub task error [%s]",
	ScpFileErrorCode:        "RequestID %s scp file %s error message %s",
	PingAgentErrorCode:      "RequestID %s ping ip %s error message %s",
	StatusMethodNotAllowedErrorCode:	"RequestID %s error message this method [%s] is not allow",
}

// ErrorMessage Show message form code.
func ErrorMessage(code StatusCode) string {
	return statusText[code]
}

//GetHTTPCode is function to map logic response code and http code
func (s *StatusCode) GetHTTPCode() int {
	switch *s {
	case ParamErrorCode, BindParamErrorCode, StatusMethodNotAllowedErrorCode:
		return http.StatusBadRequest // 客户端请求的语法错误，服务器无法理解，请求参数有误
	case DeleteDBNotExistCode:
		return http.StatusNotFound // 服务器无法根据客户端的请求找到资源
	/*case PreStateNotCorrect:
		return http.StatusUnprocessableEntity // 请求格式正确,但是由于含有语义错误,无法响应
	case NotPermissibleError:
		return http.StatusUnauthorized // 请求要求用户的身份认证
	case StatusMethodNotAllowedErrorCode:
		return http.StatusMethodNotAllowed // 客户端请求中的方法被禁止 */
	default:
		return http.StatusInternalServerError
	}
}