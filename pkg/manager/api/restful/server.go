package restful

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/common/utils/ump"
	"Infinite_train/pkg/manager/api/common"
	"Infinite_train/pkg/manager/api/request"
	"Infinite_train/pkg/manager/config"
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/satori/go.uuid"
	"github.com/tylerb/graceful"
	validatorV9 "gopkg.in/go-playground/validator.v9"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Server struct {
	*echo.Echo
	webAddr    string
	Region     string
	AdminRoles []string
	Validate   *validatorV9.Validate
}

func NewServer(config *config.Config) (*Server, error) {
	s := new(Server)
	s.webAddr = config.WebAddr
	s.Echo = echo.New()
	s.Validate = common.NewCustomValidator()
	s.Region = config.ManagerConfig.Region
	s.AdminRoles = config.ManagerConfig.AdminRoles
	return s, nil
}

func (s *Server) Run() error {
	s.RegisterContext()
	s.RegisterMiddleware()
	s.RegisterURL()
	s.Server.Addr = s.webAddr
	err := graceful.ListenAndServe(s.Server, 5*time.Second)
	return err
}

func (s *Server) RegisterURL() {
	s.GET("/v1.0/instances/:gid", UmpDecorator(s.Instances))

}

func UmpDecorator(fn func(c echo.Context) error) func(c echo.Context) error {
	return func(c echo.Context) error {
		cc := c.(*request.CustomContext)
		requestID := cc.CommonContext.RequestID
		begin := time.Now()
		err := fn(c)
		nanos := begin.UnixNano()
		millis := nanos / 1e6
		secon := nanos / 1e9
		UmpRec := &ump.Record{}
		ct := time.Unix(secon, 0)
		UmpRec.StartTime = fmt.Sprintf("%d%02d%02d%02d%02d%02d%03d", ct.Year(), ct.Month(), ct.Day(), ct.Hour(), ct.Minute(), ct.Second(), millis%1e3)
		UmpRec.AppName = "dbs_api_server"
		UmpRec.Hostname, _ = os.Hostname()
		UmpRec.Key = strings.Replace(strings.TrimPrefix(filepath.Base(handlerName(fn)), "restful.(*Server)."), `-fm`, ``, -1)
		if c.Response().Status < 300 {
			UmpRec.ProcessState = fmt.Sprintf("%d", 0)
		} else {
			UmpRec.ProcessState = fmt.Sprintf("%d", 1)
		}
		UmpRec.RequestID = requestID
		UmpRec.ElapsedTime = fmt.Sprintf("%d", int(time.Since(begin)/1000000))
		golog.Debugf(requestID, "UmpRecord:%+v", UmpRec)
		UmpRec.WriteToFile()
		return err
	}
}

func handlerName(h echo.HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

func (s *Server) RegisterMiddleware() {
	s.Use(mw.Logger())
	for _, globalSysLogger := range golog.GlobalSysLoggers {
		if globalSysLogger.Level() > golog.LevelInfo {
			continue
		}
		s.Use(mw.LoggerWithConfig(mw.LoggerConfig{
			Format: `{"time":"${time_rfc3339}","remote_ip":"${remote_ip}",` +
				`"method":"${method}","uri":"${uri}","status":${status}, "latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out}}` + "\n",
			Output: globalSysLogger,
		}))
	}
	s.Use(mw.Recover())
}

func (s *Server) RegisterContext() {
	s.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestDump, _ := httputil.DumpRequest(c.Request(), true)
			dumpStr := string(requestDump)

			var err error
			pin, err := url.QueryUnescape(c.Request().Header.Get("X-Pin"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "X-Pin unescape error !")
			}
			tenantName, err := url.QueryUnescape(c.Request().Header.Get("X-Tenant-Name"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "X-Tenant-Name unescape error !")
			}
			password, err := url.QueryUnescape(c.Request().Header.Get("X-Password"))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "X-Password unescape error !")
			}
			tenantID := c.Request().Header.Get("X-Tenant-Id")
			token := c.Request().Header.Get("X-Auth-Token")
			region := c.Request().Header.Get("X-Region")
			requestID := c.Request().Header.Get("X-Request-Id")

			if requestID == "" {
				requestID = uuid.NewV4().String()
			}
			golog.Debugf(requestID, "Request body: %+s", dumpStr)

			uri := c.Request().RequestURI
			if strings.HasSuffix(uri, "/") || strings.Contains(uri, "//") {
				return echo.NewHTTPError(http.StatusBadRequest, "url is error!")

			}

			if pin == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Header X-Pin can't be empty")
			}
			if tenantName == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Header X-Tenant-Name can't be empty")
			}
			if tenantID == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Header X-Tenant-Id can't be empty")
			}
			if password == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Header X-Password can't be empty")
			}
			if region == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "Header X-Region can't be empty")
			}
			if region != s.Region {
				msg := fmt.Sprintf("Header region:%s not equals to conf region:%s", region, s.Region)
				return echo.NewHTTPError(http.StatusBadRequest, msg)
			}
			isAdmin := false
			roles := strings.Split(c.Request().Header.Get("X-Role"), ",")
			for _, role := range roles {
				for _, adminRole := range s.AdminRoles {
					if strings.ToLower(role) == adminRole {
						isAdmin = true
						break
					}
				}
				if isAdmin {
					break
				}
			}
			baseCc := &request.CommonContext{RequestID: requestID, Token: token, TenantID: tenantID,
				TenantName: tenantName, Pin: pin, Region: region, Password: password, IsAdmin: isAdmin}
			cc := request.CustomContext{Context: c, CommonContext: baseCc}
			return h(cc)
		}
	})
}
