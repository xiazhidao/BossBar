package utils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const salt = "aaXQqWMhWlUnUrviX0eqGKPGG4qMF72b"
const Meaningless = "QqWMhW"

//将字符串加密成 md5
func String2md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	has = md5.Sum([]byte(string(has[:]) + salt))
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}

func BuildOrderNo() string {
	return time.Now().Format("20060102150405") + RandomNumString(9)
}

func GetPureNumber(origin string) string {
	var numbers []byte
	for _, v := range origin {
		if v >= 48 && v <= 57 {
			numbers = append(numbers, byte(v))
		}
	}

	return string(numbers)
}

//RandomString 在数字、大写字母、小写字母范围内生成num位的随机字符串
func RandomString(length int) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	// 小于36在大写范围内随机，其他在小写范围随机
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		t := rand.Intn(62)
		if t < 10 {
			result = append(result, strconv.Itoa(rand.Intn(10)))
		} else if t < 36 {
			result = append(result, string(rand.Intn(26)+65))
		} else {
			result = append(result, string(rand.Intn(26)+97))
		}
	}
	return strings.Join(result, "")
}

func RandomNumString(length int) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	// 小于36在大写范围内随机，其他在小写范围随机
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	for i := 0; i < length; i++ {
		result = append(result, strconv.Itoa(rand.Intn(10)))
	}
	return strings.Join(result, "")
}

func GetStringByParams(params url.Values, delKey ...string) []string {

	//获取排序前，删除不需要的key
	for _, v := range delKey {
		params.Del(v)
	}

	sorted_keys := make([]string, 0, len(params))
	for k := range params {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	sorted_values := make([]string, 0, len(params))
	for _, element := range sorted_keys {
		if len(params[element]) > 0 && len(params[element][0]) > 0 {
			param := params[element][0]
			sorted_values = append(sorted_values, param)
		}
	}
	return sorted_values
}

func GenerateSign(data string, secretKey string, h func() hash.Hash) string {
	hasher := hmac.New(h, []byte(secretKey))
	hasher.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func GenerateSortedSign(data map[string]string, secretKey string) string {
	var (
		sortedKeys []string
		signStr    string
	)
	for k, _ := range data {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, k := range sortedKeys {
		signStr += data[k]
	}
	hash := hmac.New(sha256.New, []byte(secretKey))
	hash.Write([]byte(signStr))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func ConfuseSensitiveInformation(sensitive string) string {
	return base64.StdEncoding.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(RandomNumString(3) + hex.EncodeToString([]byte(sensitive)) + RandomNumString(4) + Meaningless))))
}

func ParsingSensitiveInformation(data string) string {
	raw, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}

	raw1, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return ""
	}
	rawData := strings.Replace(string(raw1), Meaningless, "", -1)

	length := len(rawData)
	if length <= 7 {
		return ""
	}

	rawData1, err := hex.DecodeString(rawData[3 : length-4])
	if err != nil {
		return ""
	}

	return string(rawData1)
}

func CompleteImgUrl(imgStr string) string {
	if imgStr == "" {
		return ""
	}
	imgs := strings.Split(imgStr, ",")
	var completeImgs []string
	for _, v := range imgs {
		if strings.Contains(v, "certificate") {
			completeImgs = append(completeImgs, "http://acceptance-data.oss-cn-shenzhen.aliyuncs.com/"+v)
		} else {
			completeImgs = append(completeImgs, "https://funny-boys.s3.sa-east-1.amazonaws.com"+v)
		}
	}
	return strings.Join(completeImgs, ",")
}
