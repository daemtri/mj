package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"math"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func GetAuth() []rune {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var list []rune
	for i := 0; i < 6; i++ {
		ran := r.Intn(122-97+1) + 97
		list = append(list, rune(ran))
	}
	return list
}

// 验证是否邮箱
func EmailRegexp(mail string) bool {
	b := false
	if mail != "" {
		reg := regexp.MustCompile(`^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(\.[a-zA-Z0-9_-]+)$`)
		b = reg.FindString(mail) != ""
	}
	return b
}

// 验证是否手机
func PhoneRegexp(phone string) bool {
	b := false
	if phone != "" {
		reg := regexp.MustCompile(`^(86)*0*1\d{10}$`)
		b = reg.FindString(phone) != ""
	}
	return b
}

// 验证账号是否合法
func AccountRegexp(account string) bool {
	b := false
	if account != "" {
		reg := regexp.MustCompile(`^[a-zA-Z0-9]{6,8}$`)
		b = reg.FindString(account) != ""
	}
	return b
}

// 数值类型转成点分结构的IP地址
// eg: t.Log((InetTontoa(3232235966).String()))
func InetTontoa(ipnr uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// 字节类型的IP地址转成数值类型
// t.Log((InetTobton(net.IPv4(192,168,1,190))))
func InetTobton(ipnr net.IP) uint32 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// 点分结构的IP地址转成数值类型
// eg: t.Log((InetToaton("192.168.1.190")))
func InetToaton(ipnr string) uint32 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32

	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)

	return sum
}

// 验证只能由数字字母下划线组成的5-17位密码字符串
func AalidataPwd(name string) (b bool) {
	if name != "" {
		//reg := regexp.MustCompile(`^[a-zA-Z0-9_]*$`)
		reg := regexp.MustCompile(`^[a-zA-Z_]\w{5,17}$`)
		b = reg.FindString(name) != ""
	}
	return
}

// 不可见字符,用于用户提交的字符过滤分别对应为：,\0   \t  _  space  "  ` ctrl+z \n \r  `  %   \  ,
var IllegalNameRune = [13]rune{0x00, 0x09, 0x5f, 0x20, 0x22, 0x60, 0x1a, 0x0a, 0x0d, 0x27, 0x25, 0x5c, 0x2c}

var hasIllegalNameRune = func(c rune) bool {
	for _, v := range IllegalNameRune {
		if v == c {
			return true
		}
	}
	return false
}

// 限制最大字符数，检测不可见字符
// maxcount 限制的最大字符数，1个中文=2个英文
func LegalName(name string, maxcount int) bool {
	if !utf8.ValidString(name) {
		return false
	}

	num := len([]rune(name)) + len([]byte(name))
	result := float64(num) / 4.0
	sum := int(result + 0.99)

	if sum > maxcount*2 {
		return false
	}
	return strings.IndexFunc(name, hasIllegalNameRune) == -1
}

/**
 * 截取字符串
 * @param string str
 * @param begin int
 * @param length int
 * @return int 长度
 */
func SubStr(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int(x)
}

//整形转换成字节
func Int64ToBytes(n int64) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt64(b []byte) int64 {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	return int64(x)
}

//切片中字符串第一个位置
func SliceIndexOf(arr []string, str string) int {
	var index int = -1
	arrlen := len(arr)
	for i := 0; i < arrlen; i++ {
		if arr[i] == str {
			index = i
			break
		}
	}
	return index
}

//字节转换成整形
func SliceLastIndexOf(arr []string, str string) int {
	var index int = -1
	for arrlen := len(arr) - 1; arrlen > -1; arrlen-- {
		if arr[arrlen] == str {
			index = arrlen
			break
		}
	}
	return index
}

//字节转换成整形
func SliceRemoveFormSlice(oriArr []string, removeArr []string) []string {
	endArr := oriArr[:]
	for _, value := range removeArr {
		index := SliceIndexOf(endArr, value)
		if index != -1 {
			endArr = append(endArr[:index], endArr[index+1:]...)
		}
	}
	return endArr
}

// 获取本周六零点时间截
func TimestampSaturday() int64 {
	now:= time.Now()
	unix:=time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	return unix+ int64(time.Saturday-now.Weekday()) *86400
}
// 把时间戳转换成头像存储目录
func TimeToHeadphpoto(t int64, userid int, headname int64) (string, string) {
	var str, name string
	ti := time.Unix(t, 0)
	str = ti.Format("2006/01/02/15")

	str = "./headpic/" + str + "/" + strconv.Itoa(userid)
	if headname == 0 {
		name = "/130_" + strconv.Itoa(userid) + ".jpg"
	} else {
		name = "/" + strconv.Itoa(int(headname)) + ".jpg"
	}
	return str, name
}

// 把时间戳转换成头像存储目录
func TimeToPhpotoPath(t int64, userid int) string {
	var str string
	ti := time.Unix(t, 0)
	str = ti.Format("2006/01/02/15")
	return "./photo/" + str + "/" + strconv.Itoa(userid)
}
func UseridCovToInvate(userid string) uint32 {
	useridbyte := []byte(userid)
	useridbyte = useridbyte[len(useridbyte)-4:]
	timestr := []byte(strconv.Itoa(int(time.Now().Unix())))
	timestr = timestr[len(timestr)-5:]
	useridbyte = append(useridbyte, timestr...)
	code, _ := strconv.Atoi(string(useridbyte))
	return uint32(code)
}

var base = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var flipbase = flip(base)
var baselen = len(base)

func Base62encode(num uint64) string {
	baseStr := ""
	for {
		if num <= 0 {
			break
		}

		i := num % uint64(baselen)
		baseStr += base[i]
		num = (num - i) / uint64(baselen)
	}
	return baseStr
}

func Base62decode(base62 string) uint64 {
	var rs uint64 = 0
	len := uint64(len(base62))
	var i uint64
	for i = 0; i < len; i++ {
		rs += flipbase[string(base62[i])] * uint64(math.Pow(float64(baselen), float64(i)))
	}
	return rs
}

func flip(s []string) map[string]uint64 {
	f := make(map[string]uint64)
	for index, value := range s {
		f[value] = uint64(index)
	}
	return f
}

// 用gob进行数据编码
func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// 用gob进行数据解码
//
func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 对象深度拷贝
func Clone(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// 字符串加法
func StringAdd(numStr string) string {
	runeArr := []rune(numStr)
	length := len(runeArr)
	add := true
	for i := length - 1; i >= 0; i-- {
		if runeArr[i] < 57 {
			runeArr[i]++
			add = false
			break
		} else {
			runeArr[i] = 48
		}
	}
	if add {
		runeArr = append([]rune{49}, runeArr...)
	}
	return string(runeArr)
}

const FORMAT string = "2006-01-02 15:04:05"
const FORMATDATA string = "2006-01-02 "

// 获取当前时间截
func TimestampNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前时间截
func Timestamp() int64 {
	return time.Now().Unix()
}

// 获取本地当天零点时间截
func TimestampToday() int64 {
	return time.Date(Year(), Month(), Day(), 0, 0, 0, 0, time.Local).Unix()
}

// 获取本地昨天零点时间截
func TimestampYesterday() int64 {
	return TimestampToday() - 86400
}

// 获取本地明天零点时间截
func TimestampTomorrow() int64 {
	return TimestampToday() + 86400
}

// 获取当前年
func Year() int {
	return time.Now().Year()
}

// 获取当前月
func Month() time.Month {
	return time.Now().Month()
}

// 获取当前天
func Day() int {
	return time.Now().Day()
}

// 获取当前周
func Weekday() time.Weekday {
	return time.Now().Weekday()
}

// 获取指定时间截的年
func Unix2Year(t int64) int {
	return time.Unix(t, 0).Year()
}

// 获取指定时间截的月
func Unix2Month(t int64) time.Month {
	return time.Unix(t, 0).Month()
}

// 获取指定时间截的天
func Unix2Day(t int64) int {
	return time.Unix(t, 0).Day()
}

// 时间戳转str格式化时间
func Unix2Str(t int64) string {
	return time.Unix(t, 0).Format(FORMAT)
}

// str格式当前日期
func DateStr() string {
	return time.Now().Format(FORMATDATA)
}

// str格式化时间转时间戳
func Str2Unix(t string) (int64, error) {
	//the_time, err := time.Parse(FORMAT, t)
	the_time, err := time.ParseInLocation(FORMAT, t, time.Local)
	if err == nil {
		return the_time.Unix(), err
	}
	return 0, err
}

// 获取指定年月的天数
func MonthDays(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return
}

// md5 加密
// func Md5(text string) string {
// 	hashMd5 := md5.New()
// 	io.WriteString(hashMd5, text)
// 	return fmt.Sprintf("%x", hashMd5.Sum(nil))
// }
func Md5(text string) string {
	h := md5.New()
	h.Write([]byte(text))                 // 需要加密的字符串为 123456
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}

// 延迟second
func Sleep(second int) {
	<-time.After(time.Duration(second) * time.Second)
}

// 延迟1~second
func SleepRand(second int) {
	<-time.After(time.Duration(rand.Intn(second)+1) * time.Second)
}

// 延迟second
func Sleep64(second int64) {
	<-time.After(time.Duration(second) * time.Second)
}

// 延迟1~second
func SleepRand64(second int64) {
	<-time.After(time.Duration(rand.Int63n(second)+1) * time.Second)
}
