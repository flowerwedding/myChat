/**
 * @Title  handle
 * @description  #
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package server

import (
	"myChat/global"
	"myChat/logic"
	"net/http"
)

var rootDir string

func RegisterHandle(){
	global.Init()

	go logic.Broadcaster.Start()

	http.HandleFunc("/",homeHandleFunc)
	http.HandleFunc("/ws",WebSocketHandleFunc)
}