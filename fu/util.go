package fu

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
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
		Logger.Infof("getRawDataError")
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

func Keys(a map[string]string) []string {
	var b []string
	for k, _ := range a {
		b = append(b, k)
	}
	return b
}
func Values(a map[string]string) []string {
	var b []string
	for _, v := range a {
		b = append(b, v)
	}
	return b
}
func FloatToString(f float64) string {
	if f == 1.0 {
		return "1.0"
	}
	if f < 1.0 {
		//考虑到稳定性
		return fmt.Sprintf("%.5f", f)
	}
	return fmt.Sprintf("%.3f", f)
}

func UniqueMerge(a ...[]string) (newArr []string) {
	//对多个集合进行merge操作
	var ans []string
	had := map[string]bool{}
	for _, s := range a {
		if s == nil {
			continue
		}
		for _, i := range s {
			if !had[i] {
				ans = append(ans, i)
				had[i] = true
			}
		}
	}
	return ans
}
func IndexOf(a []string, v string) int {
	//获取v在字符串数组中的下标
	if a == nil {
		return -1
	}
	for ind, s := range a {
		if s == v {
			return ind
		}
	}
	return -1
}

func CopyString(a string) string {
	b := make([]byte, len(a))
	copy(b, a)
	return string(b)
}
func CopyMap(a map[string]string) map[string]string {
	b := map[string]string{}
	for k, v := range a {
		b[k] = v
	}
	return b
}

func Dict2Query(a map[string]string) string {
	//把map转化成一个固定形式的字符串，会执行url转义
	type kv struct {
		k string
		v string
	}
	var kvs []kv
	for k, v := range a {
		kvs = append(kvs, kv{k: k, v: v})
	}
	query := ""
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].k < kvs[j].k
	})
	for ind, x := range kvs {
		if ind != 0 {
			query += "&"
		}
		query += fmt.Sprintf("%v=%v", url.QueryEscape(x.k), url.QueryEscape(x.v))
	}
	return query
}
func ParseComplicateQuery(s string) map[string][]string {
	//解析形如psm=tikcast.room.permission&regions[0]=alisg&hack_name=1234的字符串，因为hertz框架不支持这种形式
	kvList := strings.Split(s, "&")
	ans := map[string][]string{}
	push := func(k, v string) {
		if ans[k] == nil {
			ans[k] = []string{v}
			return
		}
		ans[k] = append(ans[k], v)
	}
	for _, i := range kvList {
		var k, v string
		if !strings.Contains(i, "=") {
			k = i
			v = ""
		} else {
			ind := strings.Index(i, "=")
			k = i[:ind]
			v = i[ind+1:]
		}
		if strings.HasSuffix(k, "]") {
			left := strings.Index(k, "[")
			k = k[:left]
		}
		push(k, v)
	}
	return ans
}
func ArrayEqualString(m1, m2 []string) bool {
	//判断两个字符串数组是否严格相等
	if len(m1) != len(m2) {
		return false
	}
	for k, v := range m1 {
		if v != m2[k] {
			return false
		}
	}
	return true
}

func Filter(a []string, predicate func(string) bool) []string {
	//字符串过滤，只保留满足predicate的元素
	var b []string
	for _, i := range a {
		if predicate(i) {
			b = append(b, i)
		}
	}
	return b
}

func IsBlank(text string) bool {
	return len(strings.TrimSpace(text)) == 0
}
func IsBlankPtr(text *string) bool {
	return text == nil || IsBlank(*text)
}

func GetString(option *string, defaultString string) string {
	if option == nil {
		return defaultString
	}
	return *option
}
func Second2String(second int64) string {
	return time.Unix(second, 0).Format("2006-01-02 15:04:05")
}
