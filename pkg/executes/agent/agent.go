package agent

import (
	"Infinite_train/pkg/agent/config"
	"Infinite_train/pkg/agent/context"
	"Infinite_train/pkg/agent/cron"
	logCommonConfig "Infinite_train/pkg/common/config"
	"Infinite_train/pkg/common/utils"
	"Infinite_train/pkg/common/utils/linux"
	"Infinite_train/pkg/common/utils/log/golog"
	"flag"
	"fmt"
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

	var banner = "Welcome to agent!"

	var cfgFile = flag.String("config", "/etc/agent.toml", "agent configure file absolute path!")
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
	context.Agent.Config = conf
	context.Agent.LocalIP = localIP

	// init log
	err = logCommonConfig.InitConfig(conf.LogConfigs)
	if err != nil {
		golog.Error("", err.Error())
		return
	}
	golog.Info("0", banner, "time: ", time.Now())
	golog.Infof("0", "Init Config: %s", conf.String())

	// init crontab server
	cronC, err := cron.Start()
	if err != nil {
		golog.Errorx("0", "init schedule cron error: %s\n", err.Error())
		return
	}

	// init rpc server
	/*rpcServiceInit := new(rpcService.ManagerRPC)
	rpcServer, err := rpc.NewServer(conf, rpcServiceInit)
	if err != nil {
		golog.Errorx("0", "new rpc server occurs error: %s\n", err.Error())
		return
	}
	go rpcServer.Run()*/

	// SIGHUP: reload，终端控制进程结束
	// SIGINT: ctrl + c
	// SIGTERM: 结束程序(可以被捕获、阻塞或忽略)
	// SIGQUIT: 用户发送QUIT字符(Ctrl+/)触发
	// SIGPIPE: 消息管道损坏(FIFO/Socket通信时，管道未打开而进行写操作)
	exitChan := make(chan struct{})
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
				context.Agent.Config = configNew
				golog.Info("0", "SIGHUP parse config file successfully!", "signal", sig)
				golog.Infof("0", "SIGHUP Config: %s", configNew.String())
			} else if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
				golog.Info("0", "OS order me to quit, so kill myself", "signal", sig)
				for _, GlobalSysLogger := range golog.GlobalSysLoggers {
					GlobalSysLogger.Close()
				}
				//rpcServer.Close()
				close(cronC)
				golog.Close()
				close(exitChan)
				return
			} else if sig == syscall.SIGPIPE {
				golog.Info("0", "Ignore broken pipe signal", "signal", sig)
			}
		}
	}()
	<-exitChan

	if err != nil {
		golog.Errorx("0", "Restful server run occurs error: %s\n", err.Error())
		return
	}

	return
}
