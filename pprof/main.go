/**
*@author:yangzhou
*@date: 2022/9/1
*@email: yangzhou2224@shengtian.com
*@description:
 */
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main()  {
	go func() {
	//查看地址：http://127.0.0.1:6060/debug/pprof/
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	time.Sleep(time.Minute*10)
}

