package constant

const (
	InstanceStateUnauthorized			= "UNAUTHORIZED"
	InstanceStateInCheck				= "IN_CHECK"				// 正在验证账号密码、权限等信息
	InstanceStateAuthenticationFailed	= "AUTHENTICATION_FAILED"	// 账号密码错误
	InstanceStatePrivilegesNotEnough	= "PRIVILEGES_NOT_ENOUGH"	// 权限不足
	InstanceStateConnectionError		= "CONNECTION_ERROR"		// 网络连接异常
	InstanceStateAgentError				= "AGENT_ERROR"				// agent通信异常
	InstanceStateActive					= "ACTIVE"

	MetricTypeMySQL 					= "mysql"
	MetricTypeLinux						= "linux"

	MetricCollectorMethodMySQLConnect 	= "mysqlConnect"
	MetricCollectorMethodSystemAgent 	= "agent"
	MetricCollectorMethodMySQLAnsible 	= "ansible"

	InstanceNameRegEx 					= "^[A-Za-z0-9_\u4e00-\u9fa5_\\-_]{2,32}$"
)