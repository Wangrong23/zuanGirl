package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	iotqq "zuanGirl/model"
	"zuanGirl/plugins"
)

var fromUserIDStr string

func GroupMsgRouter(mess iotqq.Data) {
	fromUserIDStr = strconv.FormatInt(mess.FromUserID, 10)

	if fromUserIDStr == qq {
		return
	}

	log.Println("群聊消息: ", mess.FromNickName+"<"+fromUserIDStr+">: "+mess.Content)

	switch mess.MsgType {

	case "TextMsg":
		//文本消息
		cm := strings.Split(mess.Content, " ")
		groupTextMsg(mess, cm)

	case "PicMsg":
		//图片消息
		groupPicMsg(mess)
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
	case "AtMsg":
		//艾特某人消息
		groupAtMsg(mess)
	case "TempSessionMsg":
		//临时消息
		//群内成员向你私聊等等

	}
	return
}

func groupTextMsg(mess iotqq.Data, cm []string) {
	switch mess.Content {
	case "菜单":
		iotqq.Send(mess.FromGroupID, 2, "你好我是米娅😊\n1.赞我（50个赞哟😘）\n2.签到(正在开发)\n3.获取用户 QQ号\n4.天气 城市", 0)
		break
	case "赞我":
		ok := true
		for i := 0; i < len(zanok); i++ {
			if zanok[i] == mess.FromUserID {
				ok = false
			}
		}
		if ok {
			iotqq.Send(mess.FromGroupID, 2, "正在赞，可能需要50s时间🤣", 0)
			for i := 1; i <= 50; i++ {
				iotqq.Zan(strconv.Atoi(fromUserIDStr))
				time.Sleep(time.Second * 1)
			}
			iotqq.Send(mess.FromGroupID, 2, "已经赞了50次，如果没有成功，可能是腾讯服务器限制了！", 0)
			zanok = append(zanok, mess.FromUserID)
		} else {
			iotqq.Send(mess.FromGroupID, 2, "之前已经赞了", 0)
		}
		break
	}

	switch cm[0] {

	case "语音":
		if len(cm) < 2 {
			iotqq.Send(mess.FromGroupID, 2, "命令输入错误！", 0)
			return
		}
		iotqq.SendVoice(mess.FromGroupID, 2, cm[1])
	case "天气":
		if len(cm) < 2 {
			iotqq.Send(mess.FromGroupID, 2, "命令输入错误！", 0)
			return
		}
		n := plugins.GetWeather(cm, qq)
		iotqq.SendA(mess.FromGroupID, 2, n, "JsonMsg")
	case "获取用户":
		if len(cm) < 2 {
			iotqq.Send(mess.FromGroupID, 2, "命令输入错误！", 0)
			return
		}
		a, _ := strconv.Atoi(cm[1])
		temp := iotqq.Getinfo(a)
		var user iotqq.QQinfo
		err := json.Unmarshal([]byte(temp), &user)
		if err != nil {
			fmt.Println("反序列化出错,info:", err)
		} else {
			iotqq.Send(mess.FromGroupID, 2, "QQ昵称:"+user.Data.Nickname+"\nQQ账号:"+strconv.Itoa(user.Data.Uin)+"\nVip等级:"+strconv.Itoa(user.Data.Qqvip)+"\n绿钻等级:"+strconv.Itoa(user.Data.Greenvip)+"\n红钻等级:"+strconv.Itoa(user.Data.Redvip), 0)
		}

	}
}
func groupAtMsg(mess iotqq.Data) {
	fmt.Println(mess)

	var atContent iotqq.AtContent
	err := json.Unmarshal([]byte(mess.Content), &atContent)
	if err != nil {
		fmt.Println("反序列化出错,info:", err)
	}

	cm := strings.Split(atContent.Content, " ")

	if len(cm) > 1 {
		chatmsg := plugins.GetMsgText(cm[1])
		iotqq.Send(mess.FromGroupID, 2, " "+chatmsg, int(mess.FromUserID))
		return
	}

	//if strings.Index(atContent.Content, "喷") != -1{
	//
	//}
}

func groupPicMsg(mess iotqq.Data) {
	//FileMd5 := mess.Content
}
