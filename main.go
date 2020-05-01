package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"zuanGirl/model"
	"zuanGirl/plugins"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

var url1, qq string
var conf iotqq.Conf
var zanok, qd []int64

func init() {
	file, err := os.Open("main.conf")
	conf = iotqq.Conf{Enable: true, GData: make(map[string]int)}
	//log.Println(file)
	if err != nil {
		log.Println(err)
		os.Create("main.conf")
		f, _ := os.OpenFile("main.conf", os.O_APPEND, 0644)
		defer f.Close()
		enc := json.NewEncoder(f)
		conf.Enable = true
		conf.GData = make(map[string]int)
		enc.Encode(conf)
	}
	defer file.Close()
	tmp := json.NewDecoder(file)
	//log.Println(tmp)
	for tmp.More() {
		err := tmp.Decode(&conf)
		if err != nil {
			fmt.Println("Error:", err)
		}
		//fmt.Println(conf)
	}
}
func periodlycall(d time.Duration, f func()) {
	for x := range time.Tick(d) {
		f()
		log.Println(x)
	}
}
func resetzan() {

	m1 := len(zanok)
	for m := 0; m < m1; m++ {
		i := 0
		zanok = append(zanok[:i], zanok[i+1:]...)
	}
	m2 := len(qd)
	for m := 0; m < m2; m++ {
		i := 0
		qd = append(qd[:i], qd[i+1:]...)
	}
}
func SendJoin(c *gosocketio.Client) {
	log.Println("获取QQ号连接")
	result, err := c.Ack("GetWebConn", qq, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("emit", result)
	}
}
func save() {
	f, _ := os.OpenFile("main.conf", os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.Encode(conf)
}
func main() {
	var site string
	var port int
	port = 8888
	//fmt.Println("IOTQQ插件 - 基于SocketIO V0.0.1")
	//fmt.Println("作者:Enjoy")
	//fmt.Println("\n请输入Iotqq的Web地址(无需http://和端口): ")
	//fmt.Scan(&site)
	//fmt.Println("\n请输入Iotqq的端口号: ")
	//fmt.Scan(&port)
	//fmt.Println("\n请输入QQ机器人账号: ")
	//fmt.Scan(&qq)

	//调试用
	site = "106.54.140.137"
	qq = "2461784356"

	runtime.GOMAXPROCS(runtime.NumCPU())
	url1 = site + ":" + strconv.Itoa(port)
	iotqq.Set(url1, qq)

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(site, port, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	//收到群消息的回调事件
	err = c.On("OnGroupMsgs", func(h *gosocketio.Channel, args iotqq.Message) {
		var mess iotqq.Data = args.CurrentPacket.Data
		/*
			mess.Content 消息内容 string
			mess.FromGroupID 来源QQ群 int
			mess.FromUserID 来源QQ int64
			mess.iotqqType 消息类型 string
		*/
		log.Println("群聊消息: ", mess.FromNickName+"<"+strconv.FormatInt(mess.FromUserID, 10)+">: "+mess.Content)
		cm := strings.Split(mess.Content, " ")
		if mess.Content == "菜单" {
			iotqq.Send(mess.FromGroupID, 2, "你好我是米娅😊\n1.赞我（50个赞哟😘）\n2.签到(正在开发)\n3.获取用户 QQ号\n4.天气 城市")
			return
		}
		if mess.Content == "签到" {
			ok := true
			for i := 0; i < len(qd); i++ {
				if mess.FromUserID == 2435932516 {
					break
				}
				if qd[i] == mess.FromUserID {
					ok = false
					break
				}
			}
			if ok {
				_, err := conf.GData[strconv.FormatInt(mess.FromUserID, 10)]
				if err != false {
					conf.GData[strconv.FormatInt(mess.FromUserID, 10)] += 1
					iotqq.Send(mess.FromGroupID, 2, "签到成功 😘 当前金币:"+strconv.Itoa(conf.GData[strconv.FormatInt(mess.FromUserID, 10)]))
				} else {
					conf.GData[strconv.FormatInt(mess.FromUserID, 10)] = 1
					iotqq.Send(mess.FromGroupID, 2, "签到成功 这是你第一次签到哟😜 当前金币:"+strconv.Itoa(conf.GData[strconv.FormatInt(mess.FromUserID, 10)]))
				}
				save()
				qd = append(qd, mess.FromUserID)
			} else {
				iotqq.Send(mess.FromGroupID, 2, "已经签到过了")
			}
			return
		}
		if mess.Content == "赞我" {
			ok := true
			for i := 0; i < len(zanok); i++ {
				if zanok[i] == mess.FromUserID {
					ok = false
				}
			}
			if ok {
				iotqq.Send(mess.FromGroupID, 2, "正在赞，可能需要50s时间🤣")
				for i := 1; i <= 50; i++ {
					iotqq.Zan(strconv.Atoi(strconv.FormatInt(mess.FromUserID, 10)))
					time.Sleep(time.Second * 1)
				}
				iotqq.Send(mess.FromGroupID, 2, "已经赞了50次，如果没有成功，可能是腾讯服务器限制了！")
				zanok = append(zanok, mess.FromUserID)
			} else {
				iotqq.Send(mess.FromGroupID, 2, "之前已经赞了")
			}
			return
		}
		if cm[0] == "语音" {
			if len(cm) < 2 {
				iotqq.Send(mess.FromGroupID, 2, "命令输入错误！")
				return
			}
			iotqq.SendVoice(mess.FromGroupID, 2, cm[1])
		}
		if cm[0] == "天气" {
			if len(cm) < 2 {
				iotqq.Send(mess.FromGroupID, 2, "命令输入错误！")
				return
			}
			n := plugins.GetWeather(cm,qq)
			iotqq.SendA(mess.FromGroupID, 2, n, "JsonMsg")
		}
		if cm[0] == "获取用户" {
			if len(cm) < 2 {
				iotqq.Send(mess.FromGroupID, 2, "命令输入错误！")
				return
			}
			a, _ := strconv.Atoi(cm[1])
			temp := iotqq.Getinfo(a)
			var user iotqq.QQinfo
			err = json.Unmarshal([]byte(temp), &user)
			if err != nil {
				fmt.Println("反序列化出错,info:", err)
			} else {
				iotqq.Send(mess.FromGroupID, 2, "QQ昵称:"+user.Data.Nickname+"\nQQ账号:"+strconv.Itoa(user.Data.Uin)+"\nVip等级:"+strconv.Itoa(user.Data.Qqvip)+"\n绿钻等级:"+strconv.Itoa(user.Data.Greenvip)+"\n红钻等级:"+strconv.Itoa(user.Data.Redvip))
			}
			return
		}

		//chatmsg := plugins.GetMsgText(mess)
		//iotqq.Send(mess.FromGroupID,2,chatmsg)

	})
	if err != nil {
		log.Fatal(err)
	}

	//收到好友消息的回调事件
	err = c.On("OnFriendMsgs", func(h *gosocketio.Channel, args iotqq.Message) {
		var mess iotqq.Data = args.CurrentPacket.Data
		log.Println("私聊消息: ", mess.Content)
		if mess.Content == "菜单" {
			iotqq.Send(mess.FromUin, 1, "你好我是米娅😊\n1.赞我（50个赞哟😘）\n2.签到(正在开发)\n3.获取用户 QQ号\n4.天气 城市")
			return
		}
		chatmsg := plugins.GetMsgText(mess)
		if string(mess.FromUin) != qq {
			iotqq.Send(mess.FromUin,1,chatmsg)
		}

	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}
	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("连接成功")
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	go SendJoin(c)
	periodlycall(24*time.Hour, resetzan)
home:
	time.Sleep(600 * time.Second)
	SendJoin(c)
	goto home
	log.Println(" [x] Complete")
}
