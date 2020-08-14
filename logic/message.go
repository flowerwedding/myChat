/**
 * @Title  message
 * @description  #
 * @Author  沈来
 * @Update  2020/8/9 14:52
 **/
package logic

import (
	"fmt"
	"time"
)

type Message struct{
	User        *User            `json:"user,omitempty"`
	Type        int              `json:"type,omitempty"`
	Content     string           `json:"content,omitempty"`
	MsgTime     time.Time        `json:"msg_time,omitempty"`

	To          string           `json:"to,omitempty"`//私信
	Ats         []string         `json:"mts,omitempty"`//@人

	//Users       map[string]*User `json:"users,omitempty"`
}

const (
	MsgTypeNormal    = iota
	MsgTypeSystem
	MsgTypeUserState
	MsgTypeError
)

func NewMessage(u *User, s string) *Message {
	message := &Message{
		User:    u,
		Type:    MsgTypeNormal,
		Content: s,
		MsgTime: time.Now(),
	}

	fmt.Println("message = ",message)
	return message
}

func NewErrorMessage(s string) *Message {
	return &Message{
		User:    System,
		Type: MsgTypeError,
		Content: s,
		MsgTime: time.Now(),
	}
}

func NewWelcomeMessage(u *User) *Message{
	return &Message{
		User:    System,
		Type: MsgTypeSystem,
		Content: "Welcome to chat ," + u.Nickname + ", please remember your token :"+ u.Token,
		MsgTime: time.Now(),
	}
}

func NewNoticeMessage(message string) *Message{
	return &Message{
		User:    System,
		Type:    MsgTypeUserState,
		Content: message,
		MsgTime: time.Now(),
	}
}