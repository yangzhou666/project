/**
*@author:yangzhou
*@date: 2022/9/1
*@email: yangzhou2224@shengtian.com
*@description:
 */
package main

import (
	"fmt"
	"net/http"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

//设置websocket
//CheckOrigin防止跨站点的请求伪造
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//websocket实现
func ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close() //返回前关闭
	for {
		//读取客户端发送来到消息
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		fmt.Println("服务端收到消息:", string(message))
		fmt.Println("mt类型:", mt)
		//写入ws数据
		T := string(message) + "收到"
		msg := *(*[]byte)(unsafe.Pointer(&T))

		//服务端发送消息到客户端websocket
		err = ws.WriteMessage(mt, msg)
		if err != nil {
			break
		}
		fmt.Println("发送消息：", T)
	}
}

func main() {
	r := gin.Default()
	r.GET("/ws", ping)
	r.Run(":8080")
}
