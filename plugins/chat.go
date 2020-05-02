package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/shiguanghuxian/txai"
	"log"
)

//todo 闲聊
var APPID = "2131029764"
var APPKEY = "x4BKsoVmTHKJQvXO"

func GetMsgText(question string) string {
	txAi := txai.New(APPID, APPKEY, true)

	// 调用对应腾讯ai接口的对应函数
	val, err := txAi.NlpTextchatForText("10000", question)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return "错误"
	}
	// 打印结果
	//log.Println(err)
	js, _ := json.Marshal(val)
	log.Println(string(js))
	if val.Ret == 0 {
		return val.Data.Answer
	} else {
		return "todo 无反馈"
	}
}
