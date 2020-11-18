package history

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

func CheckType(v interface{}) {
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Float32, reflect.Float64:
		fmt.Println("float")
	case reflect.Int, reflect.Int32, reflect.Int64:
		fmt.Println("Integer")
	default:
		fmt.Println("Unknown")
	}
}

//
type Employees struct {
	EmployeeID	string
	Name		string	`format:"normal"`
	Age			int
}

func (e *Employees) UpdateAge(newAge int) {
	e.Age = newAge
}

type Customer struct {
	CookieID	string
	Name		string
	Age			int
}

func TestReflect(t *testing.T) {
	e := &Employees{
		EmployeeID: "1",
		Name:       "Mike",
		Age:        30,
	}
	fmt.Println(reflect.ValueOf(*e).FieldByName("Name"))
	//fmt.Println((reflect.ValueOf(*e)).Elem())
	if nameField, ok := reflect.TypeOf(*e).FieldByName("Name"); !ok {
		fmt.Println("not found")
	} else {
		fmt.Println(nameField.Name)
		fmt.Println(nameField.Tag.Get("format"))
	}
	reflect.ValueOf(e).MethodByName("UpdateAge").Call([]reflect.Value{reflect.ValueOf(1)})
	//reflect.ValueOf(e).MethodByName("UpdateAge").Call()
	fmt.Println(reflect.ValueOf(1))
	/*var f float64 = 12
	CheckType(f)

	fmt.Println(reflect.TypeOf(f), reflect.ValueOf(f), reflect.ValueOf(f).Type())
	*/
}

// 万能程序
// reflect.deepEqual 实现map和slice的比较
func TestDeepEqual(t *testing.T) {
	a := map[int]string{1:"one",2:"two",3:"three"}
	b := map[int]string{1:"one",2:"two",3:"three"}
	//fmt.Println(a == b)
	fmt.Println(reflect.DeepEqual(a, b))

	s1 := []int{1,2,3}
	s2 := []int{1,2,3}
	s3 := []int{1,3,2}
	fmt.Printf("s1 == s2 %t\n", reflect.DeepEqual(s1, s2))
	fmt.Printf("s1 == s3 %t\n", reflect.DeepEqual(s1, s3))

	c1 := Customer{"1", "Mike", 40}
	c2 := Customer{"1", "Mike", 40}
	fmt.Println(c1 == c2)
	fmt.Println(reflect.DeepEqual(c1, c2))
}

func fillBySetting(st interface{}, setting map[string]interface{}) error {
	// 判断传参是否是指针
	if reflect.TypeOf(st).Kind() != reflect.Ptr {
		if reflect.TypeOf(st).Elem().Kind() != reflect.Struct {
			return errors.New("first param should be a pointer to struct type")
		}
	}
	if setting != nil {
		return errors.New("setting is nil")
	}

	var (
		field	reflect.StructField
		ok		bool
	)

	for k, v := range setting {
		if field, ok = (reflect.ValueOf(st)).Elem().Type().FieldByName(k); !ok {
			continue
		}
		if field.Type == reflect.TypeOf(v) {
			vstr := reflect.ValueOf(st)
			vstr = vstr.Elem()
			vstr.FieldByName(k).Set(reflect.ValueOf(v))
		}
	}

	return nil
}

func TestFillName(t *testing.T) {
	setting := map[string]interface{}{"Name": "Mike"}
	e := Employee{}
	fillBySetting(&e, setting)
	t.Log(e)
}

type MyInt int

func TestUnsafe(t *testing.T) {
	// 不安全的
	i := 10
	t.Log(unsafe.Pointer(&i))
	f := *(*float64)(unsafe.Pointer(&i))
	t.Log(f)

	// 安全的，自定义类型转换
	a := []int{1,2,3,4,5}
	m := *(*[]MyInt)(unsafe.Pointer(&a))
	t.Log(m)

	// 原子操作，临时写到另外的内存，写完后再原子拷贝
	var shareBufPtr unsafe.Pointer
	writeDataFn := func() {
		data := []int{}
		for i:=0;i<100;i++ {
			data = append(data, i)
		}
		atomic.StorePointer(&shareBufPtr, unsafe.Pointer(&data))
	}
	readDataFn := func() {
		data := atomic.LoadPointer(&shareBufPtr)
		fmt.Println(data, *(*[]int)(data))
	}

	var wg sync.WaitGroup
	writeDataFn()
	for i:=0;i<10;i++{
		wg.Add(1)
		go func() {
			for i:=0;i<10;i++{
				writeDataFn()
				time.Sleep(time.Millisecond * 100)
			}
			wg.Done()
		}()
	}
	wg.Add(1)
	for i:=0;i<10;i++{
		wg.Add(1)
		go func() {
			for i:=0;i<10;i++{
				readDataFn()
				time.Sleep(time.Millisecond * 100)
			}
			wg.Done()
		}()
	}
}
