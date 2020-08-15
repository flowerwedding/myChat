/**
 * @Title  handle
 * @description  router
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package server

import (
	_ "myChat/docs"
	"myChat/global"
	"myChat/logic"
	"net/http"
)

var rootDir string

func RegisterHandle(){
	global.Init()

	go logic.Broadcaster.Start()

	//ginSwagger.WrapHandler(swaggerFiles.Handler)

	http.HandleFunc("/",homeHandleFunc)
	http.HandleFunc("/ws",WebSocketHandleFunc)
	http.HandleFunc("/user_list", userListHandleFunc)
}