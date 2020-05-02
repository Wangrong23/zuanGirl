package main

import (
	gosocketio "github.com/graarh/golang-socketio"
	"log"
	iotqq "zuanGirl/model"
)

func Router(c *gosocketio.Client, err error) {
	//收到群消息的回调事件
	err = c.On("OnGroupMsgs", func(h *gosocketio.Channel, args iotqq.Message) {
		var mess = args.CurrentPacket.Data
		GroupMsgRouter(mess)
	})

	if err != nil {
		log.Fatal(err)
	}

	//收到好友消息的回调事件
	err = c.On("OnFriendMsgs", func(h *gosocketio.Channel, args iotqq.Message) {
		var mess = args.CurrentPacket.Data
		FriendMsgRouter(mess)
	})
	if err != nil {
		log.Fatal(err)
	}
}
