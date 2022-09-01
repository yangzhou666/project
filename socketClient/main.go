package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type websocketClientManager struct {
	conn        *websocket.Conn
	addr        *string
	path        string
	rawQuery    string
	sendMsgChan chan string
	recvMsgChan chan string
	isAlive     bool
	timeout     int
}

// 构造函数
func NewWsClientManager(addrIp, addrPort, path, rawQuery string, timeout int) *websocketClientManager {
	addrString := addrIp + ":" + addrPort
	var sendChan = make(chan string, 10) //定义channel大小，需要及时处理消费，否则会阻塞
	var recvChan = make(chan string, 10) //定义channel大小，需要及时处理消费，否则会阻塞
	var conn *websocket.Conn
	return &websocketClientManager{
		addr:        &addrString,
		path:        path,
		conn:        conn,
		sendMsgChan: sendChan,
		recvMsgChan: recvChan,
		isAlive:     false,
		timeout:     timeout,
		rawQuery:    rawQuery,
	}
}

// 链接服务端
func (wsc *websocketClientManager) dail() error {
	var err error
	u := url.URL{Scheme: "ws", Host: *wsc.addr, Path: wsc.path, RawQuery: wsc.rawQuery}
	fmt.Println("connecting to:", u.String())
	wsc.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	wsc.isAlive = true
	log.Printf("connecting to %s 链接成功！！！", u.String())
	return nil
}

// 发送消息到服务端
func (wsc *websocketClientManager) sendMsgThread() {
	go func() {
		defer func() {
			err := recover()               // recover() 捕获panic异常，获得程序执行权。
			fmt.Println("recover()后的内容！！") // recover()后的内容会正常打印
			if err != nil {
				fmt.Println(err) // runtime error: index out of range
				wsc.isAlive = false
			}
		}()
		for {
			msg := <-wsc.sendMsgChan
			fmt.Println("发送消息:", msg)
			// websocket.TextMessage类型
			err := wsc.conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("write:", err)
				break
			}
		}
	}()
}

// 读取服务端消息
func (wsc *websocketClientManager) readMsgThread() {
	go func() {
		for {
			if wsc.conn != nil {
				_, message, err := wsc.conn.ReadMessage()
				if err != nil {
					log.Println("readErr:", err)
					wsc.isAlive = false
					// 出现错误，退出读取，尝试重连
					break
				}
				// 需要读取数据，不然会阻塞
				wsc.recvMsgChan <- string(message)

			}
		}
	}()
}

// 开启服务并重连
func (wsc *websocketClientManager) start() {
	for {
		if wsc.isAlive == false {
			if err := wsc.dail(); err == nil {
				wsc.sendMsgThread()
				wsc.readMsgThread()
				wsc.Msg()  //构造假消息
				wsc.Recv() //接收处理服务端返回到消息
			}
		}
		time.Sleep(time.Second * time.Duration(wsc.timeout))
	}
}

//模拟websocket心跳包，假数据
func (wsc *websocketClientManager) Msg() {
	go func() {
		a := 0
		for {
			wsc.sendMsgChan <- strconv.Itoa(a)
			time.Sleep(time.Second * 1)
			a += 1
		}
	}()
}

//接收处理服务端返回到消息
func (wsc *websocketClientManager) Recv() {
	go func() {
		for {
			msg, ok := <-wsc.recvMsgChan
			if ok {
				fmt.Println("收到消息：", msg)
			}
		}
	}()
}

func main() {
	defer func() {
		err := recover()               // recover() 捕获panic异常，获得程序执行权。
		fmt.Println("recover()后的内容！！") // recover()后的内容会正常打印
		if err != nil {
			fmt.Println(err) // runtime error: index out of range
		}
	}()
	wsc := NewWsClientManager("127.0.0.1", "8070", "/ws", "token=1227%7C5TjCRXftL4FpiFPQ9ajz7h4LudQCbA7oR91X0rDp", 10)
	wsc.start()

	var w1 sync.WaitGroup
	w1.Add(1)
	w1.Wait()
}
