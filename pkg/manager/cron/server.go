package cron

import (
	"Infinite_train/pkg/manager/context"
	"Infinite_train/pkg/manager/controller"
	"github.com/go-co-op/gocron"
	"math/rand"
	"sync"
	"time"
)

var scheduler = gocron.NewScheduler(time.UTC)

func registerCron() {
	scheduler.Every(context.Manager.Config.CronInterval.IntervalEveryMinute).Seconds().Do(controller.TestCon)
}

func Start() (chan struct{}, error) {
	registerCron()
	exitCh := scheduler.StartAsync()
	return exitCh, nil

	/*gocron.SetLocker(NewLocker())
	exitCh := scheduler.StartAsync()
	go Run(exitCh)
	return exitCh, nil*/
}

func Run(exitCh <-chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	restTime := 60
	timer := time.NewTimer(0)

	for {
		select {
		case <-timer.C:
			GetMonitorMutex()
			//reset the time for next trigger
			dur := restTime/2 + r.Intn(restTime)
			timer.Reset(time.Second * time.Duration(dur))

			//receive the exit signal from caller.
			//in this case, exitCh should be the returned value of gocron.Start()
		case <-exitCh:
			return
		}
	}
}

func GetMonitorMutex() {
	if GetMutex() {
		//avoid registering duplicated tasks
		if scheduler.Len() < 1 {
			registerCron()
		}
	} else {
		//clear all existed cron tasks after failing to get the mutex
		scheduler.Clear()
	}
}

func GetMutex() bool {
	/*requestID := "monitor get mutex"
	Config := context.ContextInstance.Config
	httpClient := http.NewHTTPClient(requestID)
	OwnerIp, _ := getHostIP()
	var Data struct {
		Key   string `json:"key"`
		Term  int    `json:"term"`
		Owner string `json:"owner"`
	}
	var Mutex MutexResp
	Data.Key = "monitor_lock"
	Data.Term = 180
	Data.Owner = OwnerIp
	bs, err := json.Marshal(Data)
	if err != nil {
		return false
	}
	body := bytes.NewBuffer([]byte(bs))
	result, err := httpClient.GetResponse("PUT", Config.xxx+"v1.0/mutex", body)
	if err != nil {
		return false
	}
	if result.StatusCode != 200 {
		return false
	}
	resultByte, err := ioutil.ReadAll(result.Body)
	golog.Infof(requestID, "resp body content %s", string(resultByte))
	if err != nil {
		return false
	}
	golog.Infof(requestID, "get resp body to json")
	err = json.Unmarshal(resultByte, &Mutex)
	if err != nil {
		golog.Errorf(requestID, "can not Unmarshal response body")
		return false
	}
	if Mutex.Status != "OK" {
		//golog.Errorf(requestID, Mutex.Message)
		return false
	}*/
	return true
}

// lock, https://github.com/go-co-op/gocron/blob/master/example/lock.go
type Locker struct {
	cache map[string]struct{}
	l     sync.Mutex
}

func NewLocker() *Locker {
	return &Locker{
		cache: make(map[string]struct{}),
		l:     sync.Mutex{},
	}
}

func (l *Locker) Lock(key string) (bool, error) {
	l.l.Lock()
	defer l.l.Unlock()
	if _, ok := l.cache[key]; ok {
		return false, nil
	}
	l.cache[key] = struct{}{}
	return true, nil
}

func (l *Locker) Unlock(key string) error {
	l.l.Lock()
	defer l.l.Unlock()
	delete(l.cache, key)
	return nil
}