/**
 * @Title  broadcast
 * @description  管道处理
 * @Author  沈来
 * @Update  2020/8/9 14:51
 **/
package logic

import (
	"log"
	"myChat/global"
)

type broadcaster struct{
	users map[string] *User

	enteringChannel chan *User
	leavingChannel chan *User
	messageChannel chan *Message

	// 判断该昵称用户是否可进入聊天室
	checkUserChannel chan string
	checkUserCanInChannel chan bool

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),

	enteringChannel: make(chan *User),
	leavingChannel: make(chan *User),
	messageChannel: make(chan *Message, global.MessageQueueLen),

	checkUserChannel: make(chan string),
	checkUserCanInChannel: make(chan bool),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

func (b *broadcaster) UserEntering(u *User) {
	b.enteringChannel <-u
}

func (b *broadcaster) UserLeaving(u *User) {
	b.leavingChannel <-u
}

func(b *broadcaster) Broadcast (msg *Message){
	if len(b.messageChannel) >= global.MessageQueueLen {
		log.Println("broadcast queue 满了")//channel里面满了
	}
	b.messageChannel <- msg
}

func (b *broadcaster) Start(){
	for {
		select {
		case user := <-b.enteringChannel:
			b.users[user.Nickname] = user
//			b.sendUserList()
			OfflineProcessor.Send(user)
		case user := <-b.leavingChannel:
			delete(b.users, user.Nickname)
			user.CloseMessageChannel()
//			b.sendUserList()
		case msg := <-b.messageChannel:
			if msg.To == "" {
				for _, user := range b.users{
					if user.UID == msg.User.UID {
						continue
					}
					user.MessageChannel <- msg
				}
			} else {
				if user, ok := b.users[msg.To]; ok {//私信
					user.MessageChannel <- msg
				} else {
					log.Println("user:",msg.To, "not exists!")
				}
			}
			if msg.Ats == nil {
				for _, str := range msg.Ats {//@
					if user, ok := b.users[str]; ok {//私信
						user.MessageChannel <- NewNoticeMessage("你被@了")
					} else {
						log.Println("user:",str, "not exists!")
					}
				}
			}

				OfflineProcessor.Save(msg)
		case nickname := <-b.checkUserChannel:
			if _, ok := b.users[nickname]; ok {
				b.checkUserCanInChannel <- false
			} else {
				b.checkUserCanInChannel <- true
			}
		case <-b.requestUsersChannel:
			userList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				userList = append(userList, user)
			}

			b.usersChannel <- userList
		}
	}
}

func (b *broadcaster) CanEnterRoom(nickname string) bool{
	b.checkUserChannel <- nickname

	return <-b.checkUserCanInChannel
}

func (b *broadcaster) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.usersChannel
}