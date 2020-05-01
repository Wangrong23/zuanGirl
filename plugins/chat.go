package plugins

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
	iotqq "zuanGirl/model"

	"github.com/shiguanghuxian/txai"
)

//todo 闲聊
var r *rand.Rand
var APPID = "2131029764"
var APPKEY = "x4BKsoVmTHKJQvXO"

func GetMsgText (mess iotqq.Data) string {
	txAi := txai.New(APPID, APPKEY, true)

	// 调用对应腾讯ai接口的对应函数
	val, err := txAi.NlpTextchatForText("10000",mess.Content)
	// 打印结果
	log.Println(err)
	js, _ := json.Marshal(val)
	log.Println(string(js))
	return val.Data.Answer
}

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func GetReqSign(data url.Values) string {
	str := data.Encode() + "app_key=" + APPKEY
	fmt.Println(str)
	return md5V(str)
}

func md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}