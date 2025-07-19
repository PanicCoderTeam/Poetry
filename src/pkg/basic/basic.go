package basic

import (
	"context"
	realrand "crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"net"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// 注意：本文件只包含业务无关的代码，如字符串处理等。不允许依赖工程内其他代码

// 数据库节点
const (
	ECMDBSECTION          = "ecm"
	TIMEFORMAT            = "2006-01-02 15:04:05"
	TIMEFORMAT_YYYY_MM_DD = `2006-01-02`
	TIMEFORMAT_YYYYMMDD   = `20060102`
	TIMEFORMAT_YYYYMM     = `200601`
	TIMEFORMAT_YYYY_MM    = `2006-01`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 固定盐值（实际项目应使用随机盐）
const staticSalt = "F1x3dS@lt!2023"

// 密码加密函数（盐值固定）
func HashPassword(password string) string {
	// 组合密码与盐值
	saltedPassword := password + staticSalt

	// 创建 SHA-256 哈希对象
	hasher := sha256.New()
	hasher.Write([]byte(saltedPassword))

	// 生成十六进制哈希值
	return hex.EncodeToString(hasher.Sum(nil))
}

// CheckPassword 验证密码是否匹配哈希值
func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func MakeRandUint64() uint64 {
	b := make([]byte, 8)
	_, err := realrand.Read(b)
	if err != nil {
		return rand.Uint64()
	}
	return binary.LittleEndian.Uint64(b)
}

// MakeID 生成ID
func MakeID(prefix string, count int) string {
	constChars := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y',
	}
	u := MakeRandUint64()
	result := prefix + "-"
	for i, j := 0, 0; i < count; i++ {
		if j+5 > 64 {
			j, u = 0, MakeRandUint64()
		}
		result += string(constChars[u&0x1F])
		j, u = j+5, u>>5

	}
	return result
}

/*
//MakeID 生成ID
func MakeID(prefix string, count int) string {
	constChars := []byte{
		'0', '1', '2', '3', '4', '5', '6', '7', '8',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y',
	}
	result := prefix + "-"
	for i := 0; i < count; i++ {
		result += string(constChars[rand.Intn(len(constChars))])
	}
	return result
}
*/

// MysqlLikeEscape 转义字符串中的特殊字符
func MysqlLikeEscape(str string) string {
	replacer := strings.NewReplacer(`%`, `\%`, `_`, `\_`, "\\", "\\\\")
	return replacer.Replace(str)
}

// MakeNString 生成固定格式字符串
func MakeNString(s string, sep string, count int) string {
	flag := 0
	result := ""
	for i := 0; i < count; i++ {
		if flag == 0 {
			result += s
			flag = 1
		} else {
			result += sep + s
		}
	}
	return result
}

// MakeSliceInterface 传入一个具体类型的slice，返回[]interface{}
// 注：由于有些变参函数使用args ...interface{}作为参数，调用时实参是不能使用[]string，或者[]int等类型的，需要使用[]interface{}
func MakeSliceInterface(s interface{}) ([]interface{}, error) {
	v0 := reflect.ValueOf(s)
	if v0.Kind() != reflect.Slice {
		return nil, errors.New("Not slice")
	}
	var result []interface{}
	ptr := &result
	v1 := reflect.ValueOf(ptr)
	for i := 0; i < v0.Len(); i++ {
		item := reflect.New(reflect.TypeOf(s).Elem()).Elem()
		item.Set(v0.Index(i))
		v1.Elem().Set(reflect.Append(v1.Elem(), item))
	}
	return result, nil
}

// 生成[min, max]的整数
func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// ReplaceInstanceNamePattern 实例名称里的占位符替换
// 购买多台实例，如果指定模式串{R:x}，表示生成数字[x, x+n-1]，其中n表示购买实例的数量，例如server_{R:3}，
// 购买1台时，实例显示名称为server_3；购买2台时，实例显示名称分别为server_3，server_4。
// 支持指定多个模式串{R:x}
func ReplaceInstanceNamePattern(pattern string, index int) (string, error) {
	state, count := 0, 0
	var err error
	var key, result string
	for i := 0; i < len(pattern); i++ {
		switch state {
		case 0:
			if pattern[i] == '{' {
				key, state = "", 1
				continue
			}
			result += string([]byte{pattern[i]})
		case 1:
			if pattern[i] == '{' {
				result += "{" + key
				key, state = "", 1
				if err == nil {
					err = errors.New("Bad pattern")
				}
				continue
			}
			if pattern[i] == '}' {
				state = 0
				if len(key) >= 3 && key[0:2] == "R:" {
					number, err := strconv.Atoi(key[2:])
					if err == nil {
						result += strconv.Itoa(number + index)
						count++
						continue
					}
					if err == nil {
						err = errors.New("Bad pattern")
					}
				}
				result += "{" + key + "}"
				continue
			}
			key += string([]byte{pattern[i]})
		}
	}

	if err == nil && count <= 0 {
		err = errors.New("No pattern")
	}

	return result, err
}

// MakeInstanceName 根据传入的name替换占位符，如果失败则在name后直接加数字(从1开始算)。注意index从0开始
func MakeInstanceName(name string, total, index int) string {
	if name == "" {
		return name
	}
	s, err := ReplaceInstanceNamePattern(name, index)
	if err == nil {
		return s
	}
	if total == 1 {
		return name
	}
	return name + strconv.Itoa(index+1)
}

// GetFileModTime 获取文件修改时间 返回unix时间戳
func GetFileModTime(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return -1, err
	}
	return fi.ModTime().Unix(), nil
}

// GetLocalIP 取本地ip
func GetLocalIP() string {
	ifcs, err := net.Interfaces()
	if err != nil {
		return ``
	}

	for _, ifc := range ifcs {
		if (ifc.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, _ := ifc.Addrs()
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
				continue
			}
			return ipnet.IP.String()
		}
	}
	return ``
}

//   linux:
//     'Linux机器密码需8到30位，不支持“/”置于密码首位，至少包括三项   （[a-z],[A-Z],[0-9]和[()`~!@#$%^&*-+=_|{}[]:;\'<>,.?/]的特殊符号）',
//   windows:
//     'Windows机器密码需12到30位，不支持“/”置于密码首位，至少包括三项（[a-z],[A-Z],[0-9]和[()`~!@#$%^&*-+=_|{}[]:;\'<>,.?/]的特殊符号）'

// CheckPassword 检查密码
func CheckPassword(osName string, passwd string) bool {
	name := strings.ToLower(osName)
	if !strings.Contains(name, "windows") {
		if len(passwd) < 8 || len(passwd) > 30 || passwd[0] == '/' {
			return false
		}
	} else {
		if len(passwd) < 12 || len(passwd) > 30 || passwd[0] == '/' {
			return false
		}
	}
	// 判断是否至少包含三三项
	var bools [4]bool
	var speStr = []byte{
		'(', ')', '`', '~', '!', '@', '#', '$', '%', '^', '&', '*', '-', '+', '=', '_', '|', '{', '}', '[', ']', ':',
		';', '\\', '\'', '<', '>', ',', '.', '?', '/',
	}
	for index := 0; index < len(passwd); index++ {
		if passwd[index] >= 'a' && passwd[index] <= 'z' {
			bools[0] = true
		}
		if passwd[index] >= 'A' && passwd[index] <= 'Z' {
			bools[1] = true
		}
		if passwd[index] >= '0' && passwd[index] <= '9' {
			bools[2] = true
		}
		for _, val := range speStr {
			if passwd[index] == byte(val) {
				bools[3] = true
				break
			}
		}
	}
	sum := 0
	for index := 0; index < 4; index++ {
		if bools[index] {
			sum++
		}
	}
	if sum < 3 {
		return false
	}
	return true
}

func GetStringKeySlice(m map[string]bool) []string {
	var result []string
	for key := range m {
		result = append(result, key)
	}
	return result
}

func ConvertStringPrtSlice(source []*string) []string {
	var result []string
	for _, str := range source {
		result = append(result, *str)
	}
	return result
}

func BoolToInt(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func UniqueString(s []string) []string {
	result := []string{}
	m := map[string]bool{}
	for _, v := range s {
		_, ok := m[v]
		if !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func NewPassword() string {
	n := 16
	a := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := append(a, []byte(`()~!@#$%^&*-+=_|{}[]:;\'<>,.?`+"`"+`/`)...)
	for {
		mTypes := map[int]bool{}
		k := MakeRandUint64() % uint64(len(a))
		result := []byte{a[k]}
		for i := 0; i < n-1; i++ {
			k = MakeRandUint64() % uint64(len(b))
			result = append(result, b[k])
			t := 4
			if b[k] >= 'a' && b[k] <= 'z' {
				t = 1
			} else if b[k] >= 'A' && b[k] <= 'Z' {
				t = 2
			} else if b[k] >= '0' && b[k] <= '9' {
				t = 3
			}
			mTypes[t] = true
		}
		if len(mTypes) < 3 {
			continue
		}
		return string(result)
	}
	return ""
}

func Contains(target string, s []string) bool {
	for _, i := range s {
		if i == target {
			return true
		}
	}
	return false
}

func ContainsInt64(target int64, src []int64) bool {
	for _, i := range src {
		if target == i {
			return true
		}
	}
	return false
}

func ContainsEx(target string, s []interface{}, comp func(target string, v interface{}) bool) int {
	for i, _ := range s {
		if comp(target, s[i]) {
			return i
		}
	}
	return -1
}

// 有序数组中，得到targe所在的区间序号, 采用开区间
// 0表示第一区间，len(arrary) 表示最后区间
func Range(target int64, sortedArray []int64) int {
	for i, v := range sortedArray {
		if v > target {
			return i
		}
	}
	return len(sortedArray)
}

func Min(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	} else {
		var tmp int64 = math.MaxInt64
		for _, n := range nums {
			if n < tmp {
				tmp = n
			}
		}
		return tmp
	}
}

func Max(nums ...int64) int64 {
	if len(nums) == 0 {
		return 0
	} else {
		var tmp int64 = math.MinInt64
		for _, n := range nums {
			if n > tmp {
				tmp = n
			}
		}
		return tmp
	}
}

func FMax(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	} else {
		tmp := -math.MaxFloat64
		for _, x := range nums {
			if x > tmp {
				tmp = x
			}
		}
		return tmp
	}
}

// FMaxWithIndex ...
func FMaxWithIndex(nums ...float64) (float64, int) {
	if len(nums) == 0 {
		return 0, -1
	} else {
		tmp := -math.MaxFloat64
		inx := 0
		for i, x := range nums {
			if x > tmp {
				tmp = x
				inx = i
			}
		}
		return tmp, inx
	}
}

func MinInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	} else {
		var tmp = math.MaxInt64
		for _, n := range nums {
			if n < tmp {
				tmp = n
			}
		}
		return tmp
	}
}

func MaxInt(nums ...int) int {
	if len(nums) == 0 {
		return 0
	} else {
		var tmp = math.MinInt64
		for _, n := range nums {
			if n > tmp {
				tmp = n
			}
		}
		return tmp
	}
}

func Pointer(p *string) *string {
	if len(*p) == 0 {
		return nil
	}
	return p
}

func Bitps2Mbitps(x float64) float64 {
	return x / 1000 / 1000
}

func NormalizeFilterName(filter string) string {
	var builder strings.Builder
	for _, r := range filter {
		if unicode.IsLetter(r) {
			builder.WriteRune(unicode.ToLower(r))
		}
	}
	return builder.String()
}

func InstanceFamilyAndConfigByType(instanceType string) (string, string) {
	ss := strings.Split(instanceType, ".")
	if len(ss) == 2 {
		return ss[0], ss[1]
	} else {
		return instanceType, ""
	}
}

func DeepCopyByJson(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, to)
	if err != nil {
		return err
	}

	return nil
}

// SetBySameJSONField 根据相同的json标签字段进行赋值
func SetBySameJSONField(dst interface{}, src interface{}) {
	bytes, _ := json.Marshal(src)
	json.Unmarshal(bytes, dst)
}

// InArray 字符串的index
func InArray(need string, needArr ...string) int {
	for i, v := range needArr {
		if need == v {
			return i
		}
	}
	return -1
}

func AddUpFloatSlice(dst []float64, part []float64) {
	l := Min(int64(len(dst)), int64(len(part)))
	for i := 0; i < int(l); i++ {
		dst[i] += part[i]
	}
}

func Integers2Strings(integers []int64) []string {
	result := make([]string, 0, len(integers))
	for _, i := range integers {
		result = append(result, strconv.FormatInt(i, 10))
	}
	return result
}

func MakeInterfaceSlice(src interface{}) []interface{} {
	v := reflect.ValueOf(src)
	if v.Kind() != reflect.Slice {
		return nil
	}

	result := make([]interface{}, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		result = append(result, v.Index(i).Interface())
	}
	return result
}

func ParseInt64FromStrings(src []string) []int64 {
	result := make([]int64, 0, len(src))
	for _, s := range src {
		if parseInt, err := strconv.ParseInt(s, 10, 64); err == nil {
			result = append(result, parseInt)
		}
	}
	return result
}

func ParallelCall(
	ctx context.Context,
	inputData []interface{},
	threadMaxNum int,
	callback func(context.Context, interface{}) interface{},
) []interface{} {
	chanResult := make(chan interface{}, len(inputData))
	chanthreadsLimit := make(chan byte, threadMaxNum)
	defer close(chanthreadsLimit)
	defer close(chanResult)
	for i, _ := range inputData {
		chanthreadsLimit <- '1'
		go func(ctx context.Context, data interface{}) {
			r := callback(ctx, data)
			<-chanthreadsLimit
			chanResult <- r
		}(ctx, inputData[i])
	}
	result := []interface{}{}
	for range inputData {
		r := <-chanResult
		result = append(result, r)
	}
	return result
}

func NewBoolPoint(a bool) *bool {
	return &a
}

func NewIntPoint(a int) *int {
	return &a
}

func NewInt64Point(a int64) *int64 {
	return &a
}

func NewFloat32Point(a float32) *float32 {
	return &a
}

func NewStringPoint(a string) *string {
	return &a
}

func GetFieldNameByTag(v interface{}, key, value string) string {
	v1 := reflect.TypeOf(v)
	if v1.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < v1.NumField(); i++ {
		f := v1.Field(i)
		str := strings.Split(f.Tag.Get(key), ",")[0]
		if str == value {
			return f.Name
		}
	}
	return ""
}

// Intersect 求两个字符串切片的交集
func Intersect(s0, s1 []string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)

	for _, v := range s0 {
		seen[v] = struct{}{}
	}
	for _, v := range s1 {
		if _, ok := seen[v]; ok {
			result = append(result, v)
		}
	}
	return result
}

// RemoveDup 清除字符串 slice 中的重复项，顺序会被排序
func RemoveDup(ss []string) []string {
	sort.Strings(ss)
	j := 0
	for i := 0; i < len(ss); i++ {
		if ss[i] == ss[j] {
			continue
		}
		j++
		ss[j] = ss[i]
	}
	return ss[:j+1]
}

// int32 转 []byte（大端序）
func Int32ToBytesBigEndian(n int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return buf // 示例：n=257 → [0x00 0x00 0x01 0x01]
}

// []byte 转回 int32
func BytesToInt32BigEndian(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}
