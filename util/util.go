package util

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/gin-gonic/gin"
	"github.com/samuel/go-thrift/parser"
	"github.com/weiyinfu/TwoManZeroSumGame/back/log"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"strconv"
	"strings"
)

func Assert(x bool, desc string) {
	if x {

	} else {
		panic(desc)
	}
}

func AssertEqual(x interface{}, y interface{}, desc string) {
	Assert(x == y, desc)
}

func ArrayEqual(a, b []int) bool {
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func ArrayCopy(a []int) []int {
	var b []int
	for _, v := range a {
		b = append(b, v)
	}
	return b
}
func GetIntList(s string) []int {
	//s是一个逗号隔开的int列表组成的字符串
	var a []int
	for _, i := range strings.Fields(strings.ReplaceAll(s, ",", " ")) {
		v, err := strconv.Atoi(i)
		if err != nil {
			panic(fmt.Sprintf("invalid param err=%v i=%v", err, i))
		}
		a = append(a, v)
	}
	return a
}
func IntList2String(a []int) string {
	s := make([]string, len(a))
	for ind, i := range a {
		s[ind] = strconv.Itoa(i)
	}
	return strings.Join(s, ",")
}
func HumanDuration(milliseconds int64) string {
	//把持续时间转为一个字符串
	if milliseconds == 0 {
		return "0"
	}
	mill := milliseconds % 1000
	seconds := milliseconds / 1000
	d := seconds / (3600 * 24)
	seconds %= 3600 * 24
	h := seconds / 3600
	seconds %= 3600
	m := seconds / 60
	seconds %= 60
	s := ""
	if d > 0 {
		s += fmt.Sprintf("%v天", d)
	}
	if h > 0 {
		s += fmt.Sprintf("%v小时", h)
	}
	if m > 0 {
		s += fmt.Sprintf("%v分钟", m)
	}
	if seconds > 0 {
		s += fmt.Sprintf("%v秒", seconds)
	}
	if mill > 0 {
		s += fmt.Sprintf("%v毫秒", mill)
	}
	return s
}

//整型转字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	_ = binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}
func Pause() {
	reader := bufio.NewReader(os.Stdin)
	_, _, _ = reader.ReadLine()
}
func J(x interface{}) string {
	s, _ := json.Marshal(x)
	return string(s)
}
func Bind(c *gin.Context, x interface{}) error {
	data, e := c.GetRawData()
	if e != nil {
		log.Logger.Infof("getRawDataError")
		return e
	}
	e = json.Unmarshal(data, x)
	return e
}
func CopyIntMap(ma map[int]int) map[int]int {
	var a map[int]int
	for k, v := range ma {
		a[k] = v
	}
	return a
}

func ParseThrift(filepath string) map[string]*parser.Thrift {
	p := parser.Parser{}
	file2content, thriftPath, err := p.ParseFile(filepath)
	if err != nil {
		log.Logger.Error("parse file error: file=%s,err=%s", thriftPath, err)
		panic(fmt.Sprintf("parse thrift file error=%v,filepath=%v", err, filepath))
	}
	return file2content
}
func Format(template string, ctx bson.M) string {
	//正则表达式检测所有字符串，所有的字符串都应该在ctx中出现
	s, err := raymond.Render(template, ctx)
	if err != nil {
		panic(fmt.Sprintf("render template error :template=%v,error=%v", template, err))
	}
	return s
}