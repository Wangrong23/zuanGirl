package main

import (
	"encoding/json"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"
	"zuanGirl/model"
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

	Router(c, err)

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
