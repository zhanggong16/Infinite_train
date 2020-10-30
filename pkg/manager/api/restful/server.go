package restful

import (
	"Infinite_train/pkg/manager/api/common"
	"Infinite_train/pkg/manager/config"
	"github.com/labstack/echo"
	validatorV9 "gopkg.in/go-playground/validator.v9"
)

type Server struct {
	*echo.Echo
	webAddr 	string
	Region     	string
	AdminRoles 	[]string
	Validate   	*validatorV9.Validate
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