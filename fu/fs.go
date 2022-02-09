package fu

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func MakeDir(folder string) error {
	return os.MkdirAll(folder, os.ModePerm)
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func WriteFile(filePath string, content string) int {
	cout, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Sprintf("写文件'%v'失败,err=%v", filePath, err))
	}
	byteCount, _ := io.WriteString(cout, content)
	err = cout.Close()
	if err != nil {
		panic(fmt.Sprintf("关闭文件'%v'失败", filePath))
	}
	return byteCount
}
func ReadFile(filepath string) string {
	bytesContent, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("加载文件'%v'失败", filepath))
	}
	content := string(bytesContent)
	return content
}
