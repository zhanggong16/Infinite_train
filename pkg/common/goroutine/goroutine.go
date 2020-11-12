package goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func GetGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// 得到id字符串
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}

/*func runTask(id int) string {
	time.Sleep(time.Second * time.Duration(id))
	return fmt.Sprintf("The result is from %d", id)
}*/

type CollectDataBase struct {
	Tag		string
}

type CollectSystemData struct {
	CollectDataBase
	Result	string
}

func (c *CollectSystemData)AllResponse(runTaskFunc interface{}, args interface{}) error {
	numOfRunner := 3
	timeout := 1
	ch := make(chan string, numOfRunner)
	for i := 0; i < numOfRunner; i++ {
		go func() {
			ret := runTaskFunc.(func(interface{}) string)(args)
			ch <- ret
		}()
	}

	c.Result = ""
	for j := 0; j < numOfRunner; j++ {
		select {
		case r := <-ch:
			c.Result += r + "\n"
		case <-time.After(time.Second * time.Duration(timeout)):
			c.Result += "-1" + "\n"
		}
	}
	return nil
}


/*func (cd *CollectDataBase)AllResponse(taskFunc, args interface{}, numOfRunner int, timeout time.Duration) (string, error) {
	ch := make(chan string, numOfRunner)
	finalRet := ""
	for i := 0; i < numOfRunner; i++ {
		go func() {
			ret, _ := taskFunc.(func(interface{}) (string, error))(args)

		}()
	}
	select {
	case ch<-ret:
		for j := 0; j < numOfRunner; j++ {
			finalRet += <-ch + "\n"
		}
		return finalRet, nil
	case <-time.After(time.Second * timeout):
		err := errors.New("timeout")
		return "",
	}
}
*/

