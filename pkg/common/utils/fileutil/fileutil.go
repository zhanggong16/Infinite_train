package fileutil

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//GetFileSize ...
func GetFileSize(fileName string) (size int64, isNotExist bool, err error) {
	f, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		isNotExist = true
		return
	}

	if err != nil {
		golog.Error("0", err.Error())
		return
	}
	size = f.Size()
	return
}

//MD5Sum ...
func MD5Sum(fileName string) string {
	f, err := os.Open(fileName)
	if err != nil {
		golog.Error("0", err.Error())
		return ""
	}
	defer f.Close()

	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		golog.Errorf("0", "io.Copy: %s", err.Error())
		return ""
	}
	return hex.EncodeToString(md5hash.Sum(nil))
}

//ListDir ...
func ListDir(dirPath string) (files []string, err error) {
	files = make([]string, 0, 10)
	dirFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	PathSep := string(os.PathSeparator)
	for _, fi := range dirFiles {
		if fi.IsDir() { // 忽略目录
			continue
		}
		files = append(files, dirPath+PathSep+fi.Name())
	}

	return files, nil
}

//PathExists ...
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//CreateDir ...
func CreateDir(path string) error {
	f, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	if f.IsDir() {
		return nil
	}
	return nil
}

//-------------------------------------------------------------------------------------------
//start: sort files in ascending by modtime
//--------------------------------------------------------------------------------------------

type fileInfo []os.FileInfo

func (s fileInfo) Less(i int, j int) bool {
	return s[i].ModTime().UnixNano() < s[j].ModTime().UnixNano()
}

func (s fileInfo) Len() int {
	return len(s)
}

func (s fileInfo) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

//ListDirByModTime lists files in ascending by mod time
func ListDirByModTime(dirPath string) (files []string, err error) {
	files = make([]string, 0, 10)
	dirFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	sort.Sort(fileInfo(dirFiles))

	PathSep := string(os.PathSeparator)
	for _, fi := range dirFiles {
		if fi.IsDir() { // 忽略目录
			continue
		}
		files = append(files, dirPath+PathSep+fi.Name())
	}

	return files, nil
}

//ListLogFileByModTime lists files in ascending by mod time
func ListLogFileByModTime(dirPath string) (files []string, err error) {
	files = make([]string, 0, 10)
	dirFiles, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	sort.Sort(fileInfo(dirFiles))

	for _, fi := range dirFiles {
		if fi.IsDir() {
			files = append(files, dirPath+fi.Name())
		}

	}

	return files, nil
}

//-------------------------------------------------------------------------------------------
//end: sort files in ascending by modtime
//--------------------------------------------------------------------------------------------

//IsFullStringInList ...
func IsFullStringInList(strList []string, str string) bool {
	s := "," + strings.Join(strList, ",") + ","
	if strings.Contains(s, ","+str+",") {
		return true
	}

	return false
}

//GetBaseName ...
func GetBaseName(fileName string) string {
	if len(fileName) == 0 {
		return ""
	}
	_, name := filepath.Split(fileName)
	baseName := strings.Split(name, ".")[0]
	return baseName
}

//GetBaseNameArray ...
func GetBaseNameArray(fileNames []string) []string {
	var baseName []string
	for _, file := range fileNames {
		name := GetBaseName(file)
		if name != "" {
			baseName = append(baseName, name)
		}
	}

	return baseName
}

//GetFileModTime ...
func GetFileModTime(fileName string) (modTime time.Time, isNotExist bool, err error) {
	f, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		isNotExist = true
		return
	}

	if err != nil {
		golog.Error("0", err.Error())
		return
	}
	modTime = f.ModTime()
	return
}

//MakeDir func...
func MakeDir(name string) error {
	var error error
	if !CheckFileIsExist(name) {
		error = os.Mkdir(name, os.ModePerm)
	}
	return error
}

//CheckFileIsExist ...
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//CopyFile ....
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
