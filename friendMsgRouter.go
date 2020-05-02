package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	iotqq "zuanGirl/model"
	"zuanGirl/plugins"
)

func FriendMsgRouter(mess iotqq.Data) {
	if strconv.Itoa(mess.FromUin) == qq {
		return
	}
	log.Println("私聊消息: ", mess.Content)

	switch mess.MsgType {

	case "TextMsg":
		//文本消息
		cm := strings.Split(mess.Content, " ")
		friendTextMsg(mess, cm)

	case "PicMsg":
		//图片消息
		friendPicMsg(mess)
	case "RedBagMsg":
		//红包消息
	case "BigFaceMsg":
		//大表情消息
	case "JsonMsg":
		//JSON消息
		//分享新闻消息的时候会是这个值
	case "XmltMsg":
		//XML格式消息
		//分享天气、转发多条聊天记录的时候会是这个值

	}

	return
}

func friendTextMsg(mess iotqq.Data, cm []string) {
	switch mess.Content {

	case "菜单":
		iotqq.Send(mess.FromUin, 1, "你好我是米娅😊\n1.赞我（50个赞哟😘）\n2.签到(正在开发)\n3.获取用户 QQ号\n4.天气 城市", 0)
		break
	case "喷":
		content := plugins.GetZuanResult("min")
		iotqq.Send(mess.FromUin, 1, content, 0)
		break
	case "使劲喷":
		content := plugins.GetZuanResult("high")
		iotqq.Send(mess.FromUin, 1, content, 0)
		break
	default:
		chatmsg := plugins.GetMsgText(mess.Content)
		iotqq.Send(mess.FromUin, 1, chatmsg, 0)

	}
}

func friendPicMsg(mess iotqq.Data) {
	var picContent iotqq.PicContent
	err := json.Unmarshal([]byte(mess.Content), &picContent)
	if err != nil {
		fmt.Println("反序列化出错,info:", err)
	}
	fmt.Println(picContent.FileMd5)

	switch picContent.FileMd5 {

	case "tAf3CKLGpQY0IJjffKxKVw==":
		//色图来
		break
	case "uIuiLSSbPVAv+Zn6x54FeA==":
		//不够色
		break

	}
}
