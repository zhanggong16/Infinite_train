package manager

import (
	logCommonConfig "Infinite_train/pkg/common/config"
	"Infinite_train/pkg/common/utils"
	"Infinite_train/pkg/common/utils/linux"
	"Infinite_train/pkg/common/utils/log/golog"
	"Infinite_train/pkg/common/utils/mysql"
	"Infinite_train/pkg/manager/api/restful"
	"Infinite_train/pkg/manager/config"
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/model/bean"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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
	localIP, err := linux.GetLocalIP()
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
	bean.DbEngine, err = mysql.CreateOrmEngine(conf.DataBase.Account, conf.DataBase.Password, conf.DataBase.IP,
		conf.DataBase.Port, conf.DataBase.Schema, conf.DataBase.Charset, conf.DataBase.MaxIdle, conf.DataBase.MaxOpen)
	if err != nil {
		golog.Errorx("0", "Connect db error:%v\n", err.Error())
		return
	}
	bean.DbEngine.ShowSQL(false)
	golog.Infof("0", "Init db client successfully!")

	// init interface
	/*controller.InitControllerLayer()
	service.InitServiceLayer()*/

	// init restful server
	server, err := restful.NewServer(conf)
	if err != nil {
		golog.Errorx("0", "New restful server occurs error: %s\n", err.Error())
		return
	}

	// SIGHUP: reload，终端控制进程结束
	// SIGINT: ctrl + c
	// SIGTERM: 结束程序(可以被捕获、阻塞或忽略)
	// SIGQUIT: 用户发送QUIT字符(Ctrl+/)触发
	// SIGPIPE: 消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE)
	go func() {
		for {
			sig := <-sigs
			if sig == syscall.SIGHUP {
				configNew, err := config.ParseConfig(*cfgFile)
				if err != nil {
					golog.Errorf("0", "Parse config file failed!", "signal", sig)
					return
				}
				context.Instance.Config = configNew
				golog.Info("0", "SIGHUP parse config file successfully!", "signal", sig)
				golog.Infof("0", "SIGHUP Config: %s", configNew.String())
			} else if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				golog.Info("0", "OS order me to quit, so kill myself", "signal", sig)
				server.Close()
				//rpc.Close()
				for _, GlobalSysLogger := range golog.GlobalSysLoggers {
					GlobalSysLogger.Close()
				}
				golog.Close()
				return
			} else if sig == syscall.SIGPIPE {
				golog.Info("0", "Ignore broken pipe signal", "signal", sig)
			}
		}
	}()

	err = server.Run()
	if err != nil {
		golog.Errorx("0", "Restful server run occurs error: %s\n", err.Error())
		return
	}

	return
}
