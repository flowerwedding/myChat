/**
 * @Title  user
 * @description  #
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package logic

import (
	"context"
	"errors"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"sync/atomic"
	"time"
)

var globalUID uint32 = 0

type User struct {
	UID         int               `json:"uid"`
	Nickname    string            `json:"nickname"`
	EnterAt     time.Time         `json:"enter_at"`
	Addr        string            `json:"addr"`
	MessageChannel  chan *Message `json:"-"`

	conn *websocket.Conn
}

var System = &User{UID: -1,Nickname:"System"}
//var System *User

func NewUser(conn *websocket.Conn, nickname string, addr string) *User {
	user := &User{
		Nickname:       nickname,
		Addr:           addr,
		EnterAt:        time.Now(),
		MessageChannel: make(chan *Message, 32),

		conn: conn,
	}

	if user.UID == 0 {
		user.UID = int(atomic.AddUint32(&globalUID, 1))
	}

	return user
}

func (u *User) SendMessage(ctx context.Context){
	for msg := range u.MessageChannel {
		txt := msg.User.Nickname + "(" + msg.MsgTime.Format("15:04:05" ) + ")" + ":" + msg.Content
		_ = wsjson.Write(ctx, u.conn, txt)
	}
}

func (u *User) ReceiveMessage(ctx context.Context) error{
	var (
	//	receiveStr string
		receiveMsg map[string]string
		err        error
	)

	for {
	//	err = wsjson.Read(ctx, u.conn, &receiveStr)//输入字符串 "XXX"
		err = wsjson.Read(ctx, u.conn, &receiveMsg)//输入JSON格式 {"content":"XXX"}
		if err != nil {
			var closeErr websocket.CloseError
			if errors.As(err, &closeErr) {
				return nil
			}

			return err
		}

	//	sendMsg := NewMessage(u, receiveStr)
		sendMsg := NewMessage(u, receiveMsg["content"])
		Broadcaster.Broadcast(sendMsg)
	}
}

func (u *User) CloseMessageChannel(){
	close(u.MessageChannel)
}