/**
 * @Title  websocket
 * @description  服务端总流程
 * @Author  沈来
 * @Update  2020/8/9 14:53
 **/
package server

import (
	"log"
	"myChat/logic"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// @Summary  发送消息
// @Produce  json
// @Param  nickname body string true "nickname"
// @Param  token body string false "token"
// @Success  200 {object} logic.Message "成功"
// @Router  /ws [GET]
func WebSocketHandleFunc(w http.ResponseWriter,req *http.Request) {
	conn, err := websocket.Accept(w, req,&websocket.AcceptOptions{InsecureSkipVerify:true})
	if err != nil {
		log.Println("websocket accept error:", err)
		return
	}

	//新用户进来，创建该用户的实例
	//nickname := req.FormValue("nickname")
	vars := req.URL.Query()
	nickname := vars["nickname"][0]
	token := vars["token"][0]
	if l := len(nickname); l < 2 || l > 20 {
		log.Println("nickname illegal: ",nickname)
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("非法昵称，昵称长度：4-20"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname illegal!")
		return
	}
	if logic.Broadcaster.CanEnterRoom(nickname) {
		log.Println("昵称已经存在：",nickname)
		_ = wsjson.Write(req.Context(), conn, logic.NewErrorMessage("该昵称已经存在！"))
		_ = conn.Close(websocket.StatusUnsupportedData, "nickname exists!")
		return
	}

	userHasToken := logic.NewUser(conn,token, nickname, req.RemoteAddr)

	//向所有的用户告知新用户的到来
	msg := logic.NewNoticeMessage(nickname + "加入了聊天室")
	logic.Broadcaster.Broadcast(msg)

	//开启给用户发送消息的 goroutine
	go userHasToken.SendMessage(req.Context())

	//给新用户发送欢迎消息
	logic.Sending(nil,logic.NewWelcomeMessage(userHasToken),"*."+userHasToken.Nickname+".*")

	logic.Sending(nil,logic.NewWelcomeMessage(userHasToken),"*."+userHasToken.Nickname+".*")

	// 避免 token 泄露
	tmpUser := *userHasToken
	user := &tmpUser
	user.Token = ""

	//将该用户加入广播器的用户列表
	logic.Broadcaster.UserEntering(user)
	log.Println("user:", nickname, "joins chat")
/*
	go func() {
		d := 5 * time.Minute
		timer := time.NewTimer(d)
		for {
			<-timer.C
			if user.Active {
				timer.Reset(d)
				user.Active = false
			}else {
				conn.Close(websocket.StatusInternalError, "")
			}
		}
	}()*/

	//接收用户消息
	err = user.ReceiveMessage(req.Context())

	//用户离开
	logic.Broadcaster.UserLeaving(user)
	msg = logic.NewNoticeMessage(user.Nickname + "离开了聊天室")
	logic.Broadcaster.Broadcast(msg)
	log.Println("user:", nickname, "leaves chat")

	//根据读取时的错误执行不同的close
	if err != nil {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	} else {
		log.Println("read from client error:", err)
		_ = conn.Close(websocket.StatusInternalError, "Read from client error")
	}
}