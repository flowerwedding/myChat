/**
 * @Title  broadcast
 * @description  管道处理
 * @Author  沈来
 * @Update  2020/8/9 14:51
 **/
package logic

import "sync"

type broadcaster struct{
	users map[string] *User
}

var Broadcaster = &broadcaster{
	users: make(map[string]*User),
}

func (b *broadcaster) UserEntering(u *User) {
	Sending(u,nil,"entering.*.*")
}

func (b *broadcaster) UserLeaving(u *User) {
	Sending(u,nil,"*.*.leaving")
}

func(b *broadcaster) Broadcast (msg *Message){
	Sending(nil,msg,"*.message.*")
}

func (b *broadcaster) Start(){
	go b.ReceiveUser("entering.*.*")
	go b.ReceiveUser("*.*.leaving")
	go b.ReceiveMessage()
}

var lock sync.Mutex

//昵称是否重复，是否可进入
func (b *broadcaster) CanEnterRoom(nickname string) bool{
	lock.Lock()
	_, ok := b.users[nickname]
	lock.Unlock()
	return ok
}

//获取用户列表
func (b *broadcaster) GetUserList() []*User {
	lock.Lock()
	userList := make([]*User, 0, len(b.users))
	for _, user := range b.users {
		userList = append(userList, user)
	}
	lock.Unlock()
	return userList
}