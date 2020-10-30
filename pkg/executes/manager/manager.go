package manager

import (
	logCommonConfig "Infinite_train/pkg/common/config"
	"Infinite_train/pkg/common/utils"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/manager/config"
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/model/bean"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"time"
)

func Main(versionInfo *utils.VersionInfo) {

	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 1<<20)
			stack = stack[:runtime.Stack(stack, true)]
			golog.Error("0", "manager panic, err: %s, stack: %s", err, stack)
		}
	}()

	var banner = "Welcome to infinite train manager!"

	var cfgFile = flag.String("config", "/etc/manager.toml", "manager configure file absolute path!")
	var isShowVersion = false
	flag.BoolVar(&isShowVersion, "version", false, "Show version")
	flag.Parse()

	if isShowVersion {
		utils.ShowVersion(versionInfo, banner)
		return
	}

	// init config file
	if len(*cfgFile) <= 0 {
		fmt.Printf("The absolute path of the config file lost!\n")
		return
	}
	conf, err := config.ParseConfig(*cfgFile)
	if err != nil {
		fmt.Printf("Parse config file failed!\n")
		return
	}
	localIP, err := utils.GetLocalIP()
	if err != nil {
		fmt.Printf("Get local IP failed!\n")
		return
	}
	context.Instance.Config = conf
	context.Instance.ManagerIP = localIP

	// init log
	err = logCommonConfig.InitConfig(conf.LogConfigs)
	if err != nil {
		golog.Error("", err.Error())
		return
	}
	golog.Info("0", banner, "time: ", time.Now())
	golog.Infof("0", "Init Config: %s", conf.String())

	// init meta db
	/*iv := conf.ProductLineValue
	conf.DataBase.Account = encryption.Decrypt(conf.DataBase.Account, iv)
	conf.DataBase.Password = encryption.Decrypt(conf.DataBase.Password, iv)
	conf.DataBase.Schema = encryption.Decrypt(conf.DataBase.Schema, iv)*/
	bean.DbEngine, err = utils.CreateOrmEngine(conf.DataBase.Account, conf.DataBase.Password, conf.DataBase.IP,
		conf.DataBase.Port, conf.DataBase.Schema, conf.DataBase.Charset, conf.DataBase.MaxIdle, conf.DataBase.MaxOpen)
	if err != nil {
		golog.Errorx("0", "connect db error:%v\n", err.Error())
		return
	}
	bean.DbEngine.ShowSQL(false)
	golog.Infof("0", "Init db client successfully!")

	time.Sleep(2 * time.Second)

	return
}
